/**
 * NIIMBOT B21 Label Printer — Web Bluetooth Integration
 *
 * Protocol reverse-engineered from:
 *   - NiimBlueLib (github.com/MultiMote/niimbluelib)
 *   - niimprint  (github.com/AndBondStyle/niimprint)
 *
 * Supports B21, B1, D11, D110 over Web Bluetooth.
 *
 * Usage:
 *   const printer = new NiimbotPrinter();
 *   await printer.connect();
 *   await printer.printLabel(canvas, { density: 3, labelType: 1 });
 *   printer.disconnect();
 */

// ============================================================================
// CONSTANTS
// ============================================================================

const BLE_SERVICE_UUID = "e7810a71-73ae-499d-8c15-faa9aef0c3f2";

const HEAD = [0x55, 0x55];
const TAIL = [0xaa, 0xaa];

/** Request command IDs (client → printer) */
const CMD = {
  Connect:              0xc1,
  Heartbeat:            0xdc,
  SetDensity:           0x21,
  SetLabelType:         0x23,
  PrintStart:           0x01,
  PageStart:            0x03,
  SetPageSize:          0x13,
  PrintBitmapRow:       0x85,
  PrintBitmapRowIndexed:0x83,
  PrintEmptyRow:        0x84,
  PrinterCheckLine:     0x86,
  PageEnd:              0xe3,
  PrintEnd:             0xf3,
  PrinterInfo:          0x40,
  PrintStatus:          0xa3,
};

/** Response command IDs (printer → client) */
const RESP = {
  In_Connect:     0xc2,
  In_SetDensity:  0x31,
  In_SetLabelType:0x33,
  In_PrintStart:  0x02,
  In_PageStart:   0x04,
  In_SetPageSize: 0x14,
  In_PageEnd:     0xe4,
  In_PrintEnd:    0xf4,
  In_Heartbeat:   0xde,
  In_HeartbeatAdv1: 0xdd,
  In_HeartbeatAdv2: 0xd9,
  In_PrintError:  0xdb,
  In_NotSupported:0x00,
  In_PrinterCheckLine: 0xd3,
};

// ============================================================================
// PACKET CONSTRUCTION
// ============================================================================

/**
 * Build a standard NIIMBOT packet.
 * Format: [0x55, 0x55, CMD, LEN, ...DATA, CHECKSUM, 0xAA, 0xAA]
 * Checksum = XOR(CMD, LEN, ...DATA)
 */
function buildPacket(cmd, data = []) {
  const len = data.length;
  let checksum = cmd ^ len;
  for (const b of data) checksum ^= b;

  return new Uint8Array([...HEAD, cmd, len, ...data, checksum, ...TAIL]);
}

/**
 * Build a connect packet (prepends 0x03).
 */
function buildConnectPacket() {
  const inner = buildPacket(CMD.Connect, [0x01]);
  const out = new Uint8Array(1 + inner.length);
  out[0] = 0x03;
  out.set(inner, 1);
  return out;
}

/**
 * Encode a 16-bit unsigned integer as big-endian [hi, lo].
 */
function u16be(val) {
  return [(val >> 8) & 0xff, val & 0xff];
}

// ============================================================================
// IMAGE ENCODING
// ============================================================================

/**
 * Encode a canvas element into NIIMBOT row-based bitmap format.
 *
 * The B21 prints left-to-right, so the image is rotated 90° clockwise:
 *   - Canvas width  → print height (rows)
 *   - Canvas height → print width  (columns, must be multiple of 8)
 *
 * Each row is packed into bytes (8 pixels per byte, MSB = leftmost).
 * Non-white pixels are treated as black (thermal print).
 *
 * @param {HTMLCanvasElement} canvas - Source image
 * @returns {{ rows: number, cols: number, rowsData: Array }}
 */
function encodeImage(canvas) {
  const ctx = canvas.getContext("2d");
  const w = canvas.width;
  const h = canvas.height;
  const imageData = ctx.getImageData(0, 0, w, h);
  const pixels = imageData.data; // RGBA

  // After 90° CW rotation: rows = w, cols = h (rounded up to multiple of 8)
  const cols = Math.ceil(h / 8) * 8;
  const rows = w;
  const bytesPerRow = cols / 8;

  const rowsData = [];
  let prevRowData = null;
  let repeatCount = 0;

  for (let outRow = 0; outRow < rows; outRow++) {
    // Source pixel for rotated position:
    // 90° CW: src(x,y) → dst(y, width-1-x)
    // So for output row r, col c: srcX = width-1-r, srcY = c
    const srcX = w - 1 - outRow;

    const rowBytes = new Uint8Array(bytesPerRow);
    let blackCount = 0;

    for (let outCol = 0; outCol < cols; outCol++) {
      const srcY = outCol;
      let isBlack = false;

      if (srcY < h) {
        const idx = (srcY * w + srcX) * 4;
        const r = pixels[idx], g = pixels[idx + 1], b = pixels[idx + 2], a = pixels[idx + 3];
        // Non-white and non-transparent = black
        isBlack = a > 128 && (r < 200 || g < 200 || b < 200);
      }

      if (isBlack) {
        rowBytes[Math.floor(outCol / 8)] |= (0x80 >> (outCol % 8));
        blackCount++;
      }
    }

    // Check if row is identical to previous for compression
    const same = prevRowData && arraysEqual(rowBytes, prevRowData);

    if (same) {
      repeatCount++;
    } else {
      // Flush previous row if any
      if (prevRowData !== null) {
        rowsData.push({
          type: blackCount === 0 && !prevHadBlack ? "void" : "pixels",
          data: prevRowData,
          repeat: repeatCount,
          blackCount: prevBlackCount,
        });
      }
      prevRowData = rowBytes;
      prevBlackCount = blackCount;
      prevHadBlack = blackCount > 0;
      repeatCount = 1;
    }

    // Insert check line every 200 rows
    if (outRow > 0 && outRow % 200 === 0) {
      // Flush current
      if (prevRowData !== null) {
        rowsData.push({
          type: "pixels",
          data: prevRowData,
          repeat: repeatCount,
          blackCount: prevBlackCount,
        });
        prevRowData = null;
        repeatCount = 0;
      }
      rowsData.push({ type: "check", data: null, repeat: 0, blackCount: 0 });
    }
  }

  // Flush last row
  if (prevRowData !== null) {
    rowsData.push({
      type: prevBlackCount === 0 ? "void" : "pixels",
      data: prevRowData,
      repeat: repeatCount,
      blackCount: prevBlackCount,
    });
  }

  return { rows, cols, rowsData };
}

// Temp variables used during encoding (avoid closures)
let prevBlackCount = 0;
let prevHadBlack = false;

function arraysEqual(a, b) {
  if (a.length !== b.length) return false;
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false;
  }
  return true;
}

/**
 * Count total black pixels in packed row data for the pixel count header.
 */
function countBlackPixels(data, printheadPixels) {
  let count = 0;
  const limit = Math.ceil(printheadPixels / 8);
  for (let i = 0; i < Math.min(data.length, limit); i++) {
    let byte = data[i];
    while (byte) {
      count += byte & 1;
      byte >>= 1;
    }
  }
  return count;
}

// ============================================================================
// PRINTER CLASS
// ============================================================================

class NiimbotPrinter {
  constructor() {
    this.device = null;
    this.channel = null;
    this.connected = false;
    this.responseResolvers = new Map();
    this.packetBuf = [];
    this._onStatus = null;
  }

  /**
   * Connect to a NIIMBOT printer via Web Bluetooth.
   * Opens the browser's device picker filtered to NIIMBOT devices.
   */
  async connect() {
    if (!navigator.bluetooth) {
      throw new Error("Web Bluetooth not supported. Use Chrome/Edge on desktop or Android.");
    }

    this._status("Scanning for NIIMBOT printer...");

    this.device = await navigator.bluetooth.requestDevice({
      filters: [
        { namePrefix: "D11" },
        { namePrefix: "D110" },
        { namePrefix: "B1_" },
        { namePrefix: "B21_" },
        { namePrefix: "B21" },
        { namePrefix: "B1" },
        { namePrefix: "NIIMBOT" },
      ],
      optionalServices: [BLE_SERVICE_UUID],
    });

    this._status("Connecting to " + this.device.name + "...");

    this.device.addEventListener("gattserverdisconnected", () => {
      this.connected = false;
      this._status("Printer disconnected");
    });

    const server = await this.device.gatt.connect();
    const service = await server.getPrimaryService(BLE_SERVICE_UUID);
    const chars = await service.getCharacteristics();

    // Find characteristic with notify + writeWithoutResponse
    for (const ch of chars) {
      if (ch.properties.notify && ch.properties.writeWithoutResponse) {
        this.channel = ch;
        break;
      }
    }

    if (!this.channel) {
      throw new Error("No suitable BLE characteristic found on this printer");
    }

    // Start notifications for responses
    await this.channel.startNotifications();
    this.channel.addEventListener("characteristicvaluechanged", (e) => {
      this._onData(new Uint8Array(e.target.value.buffer));
    });

    // Send connect command
    await this._sendRaw(buildConnectPacket());
    await this._waitResponse(RESP.In_Connect, 3000);

    this.connected = true;
    this._status("Connected to " + this.device.name);
  }

  /**
   * Disconnect from the printer.
   */
  disconnect() {
    if (this.device && this.device.gatt.connected) {
      this.device.gatt.disconnect();
    }
    this.connected = false;
    this.device = null;
    this.channel = null;
  }

  /**
   * Print an image from a canvas element.
   *
   * @param {HTMLCanvasElement} canvas - The label image to print
   * @param {Object} opts - Print options
   * @param {number} opts.density - Print density 1-5 (default 3)
   * @param {number} opts.labelType - Label type: 1=die-cut, 2=gap, 3=transparent (default 1)
   * @param {number} opts.quantity - Number of copies (default 1)
   */
  async printLabel(canvas, opts = {}) {
    const density = opts.density || 3;
    const labelType = opts.labelType || 1;
    const quantity = opts.quantity || 1;

    if (!this.connected) {
      throw new Error("Printer not connected");
    }

    this._status("Encoding image...");
    const image = encodeImage(canvas);

    // Printhead width for B21 is typically 384 pixels (48mm * 8px/mm at 203 DPI)
    const printheadPixels = 384;

    this._status("Setting density " + density + "...");
    await this._sendAndWait(CMD.SetDensity, [density], RESP.In_SetDensity);

    this._status("Setting label type " + labelType + "...");
    await this._sendAndWait(CMD.SetLabelType, [labelType], RESP.In_SetLabelType);

    this._status("Starting print job...");
    await this._sendAndWait(CMD.PrintStart, [0x01], RESP.In_PrintStart);

    for (let copy = 0; copy < quantity; copy++) {
      this._status("Printing" + (quantity > 1 ? " copy " + (copy + 1) + "/" + quantity : "") + "...");

      await this._sendAndWait(CMD.PageStart, [], RESP.In_PageStart);
      await this._sendAndWait(CMD.SetPageSize, [
        ...u16be(image.rows),
        ...u16be(image.cols),
      ], RESP.In_SetPageSize);

      // Send image data row by row
      let pos = 0;
      for (const row of image.rowsData) {
        if (row.type === "check") {
          await this._send(CMD.PrinterCheckLine, [...u16be(pos)]);
          // Wait briefly for printer to confirm
          await this._waitResponse(RESP.In_PrinterCheckLine, 2000).catch(() => {});
          continue;
        }

        if (row.type === "void") {
          // Empty rows — skip with position advance
          pos += row.repeat;
          continue;
        }

        // Pixel row — send bitmap data
        const blackCount = countBlackPixels(row.data, printheadPixels);
        const countParts = u16be(blackCount);

        const payload = [
          ...u16be(pos),
          ...countParts,
          ...u16be(row.repeat),
          ...row.data,
        ];

        await this._send(CMD.PrintBitmapRow, payload);

        pos += row.repeat;

        // Small delay to avoid overwhelming the printer buffer
        if (pos % 50 === 0) {
          await sleep(10);
        }
      }

      await this._sendAndWait(CMD.PageEnd, [], RESP.In_PageEnd, 5000);
    }

    this._status("Finishing print job...");
    await this._sendAndWait(CMD.PrintEnd, [0x01], RESP.In_PrintEnd, 10000);

    this._status("Print complete!");
  }

  /**
   * Set a callback for status messages.
   * @param {function(string)} fn
   */
  onStatus(fn) {
    this._onStatus = fn;
  }

  // --- Internal methods ---

  _status(msg) {
    if (this._onStatus) this._onStatus(msg);
  }

  async _sendRaw(bytes) {
    if (!this.channel) throw new Error("Not connected");
    await this.channel.writeValueWithoutResponse(bytes.buffer);
  }

  async _send(cmd, data = []) {
    await this._sendRaw(buildPacket(cmd, data));
    await sleep(5); // Inter-packet delay
  }

  async _sendAndWait(cmd, data, expectedResp, timeout = 2000) {
    await this._send(cmd, data);
    return this._waitResponse(expectedResp, timeout);
  }

  _waitResponse(cmdId, timeout = 2000) {
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        this.responseResolvers.delete(cmdId);
        reject(new Error("Timeout waiting for response 0x" + cmdId.toString(16)));
      }, timeout);

      this.responseResolvers.set(cmdId, (data) => {
        clearTimeout(timer);
        resolve(data);
      });
    });
  }

  _onData(bytes) {
    // Accumulate bytes and parse packets
    for (const b of bytes) {
      this.packetBuf.push(b);
    }
    this._parsePackets();
  }

  _parsePackets() {
    while (this.packetBuf.length >= 6) {
      // Find packet header
      const headIdx = this._findHead();
      if (headIdx < 0) {
        this.packetBuf = [];
        return;
      }
      if (headIdx > 0) {
        this.packetBuf.splice(0, headIdx);
      }

      // Check we have enough bytes for header + cmd + len
      if (this.packetBuf.length < 4) return;

      const cmd = this.packetBuf[2];
      const dataLen = this.packetBuf[3];
      const totalLen = 2 + 1 + 1 + dataLen + 1 + 2; // HEAD + CMD + LEN + DATA + CHK + TAIL

      if (this.packetBuf.length < totalLen) return;

      // Extract and verify tail
      if (this.packetBuf[totalLen - 2] !== 0xaa || this.packetBuf[totalLen - 1] !== 0xaa) {
        // Bad tail, skip this head and try again
        this.packetBuf.splice(0, 2);
        continue;
      }

      const data = this.packetBuf.slice(4, 4 + dataLen);
      this.packetBuf.splice(0, totalLen);

      // Resolve any waiting promise
      const resolver = this.responseResolvers.get(cmd);
      if (resolver) {
        this.responseResolvers.delete(cmd);
        resolver(data);
      }

      // Check for errors
      if (cmd === RESP.In_PrintError) {
        this._status("Printer error: code " + (data[0] || "unknown"));
      }
    }
  }

  _findHead() {
    for (let i = 0; i < this.packetBuf.length - 1; i++) {
      if (this.packetBuf[i] === 0x55 && this.packetBuf[i + 1] === 0x55) {
        return i;
      }
    }
    return -1;
  }
}

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

// ============================================================================
// LABEL RENDERER
// ============================================================================

/**
 * Render a label canvas with QR code and text for the NIIMBOT B21.
 *
 * Default layout for 40x30mm label at 203 DPI (320x240 pixels):
 *
 *   ┌──────────────────────────────────┐
 *   │            TITLE TEXT            │
 *   │  ┌──────────┐                   │
 *   │  │          │  tag1, tag2       │
 *   │  │  QR CODE │                   │
 *   │  │          │  2024-01-15       │
 *   │  └──────────┘                   │
 *   │       pochita.synology.me/...   │
 *   └──────────────────────────────────┘
 *
 * @param {Object} opts
 * @param {string} opts.qrImageUrl - URL of the QR code PNG from server
 * @param {string} opts.title      - Note title
 * @param {string} opts.url        - Share URL text
 * @param {string} opts.date       - Date string
 * @param {string[]} opts.tags     - Tag list
 * @param {number} opts.width      - Label width in pixels (default 320)
 * @param {number} opts.height     - Label height in pixels (default 240)
 * @returns {Promise<HTMLCanvasElement>}
 */
async function renderLabel(opts) {
  const width = opts.width || 320;
  const height = opts.height || 240;

  const canvas = document.createElement("canvas");
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext("2d");

  // White background
  ctx.fillStyle = "#ffffff";
  ctx.fillRect(0, 0, width, height);

  // Load QR code image
  const qrImg = await loadImage(opts.qrImageUrl);

  // Layout calculations
  const padding = 8;
  const qrSize = Math.min(height - 60, width * 0.45); // ~45% of width
  const qrX = padding;
  const qrY = 30;

  // Title
  ctx.fillStyle = "#000000";
  ctx.font = "bold 18px sans-serif";
  ctx.textBaseline = "top";
  const titleText = truncate(opts.title || "Untitled", 30);
  ctx.fillText(titleText, padding, padding);

  // QR Code
  ctx.drawImage(qrImg, qrX, qrY, qrSize, qrSize);

  // Right side info
  const infoX = qrX + qrSize + 10;
  const infoW = width - infoX - padding;

  // Tags
  if (opts.tags && opts.tags.length > 0) {
    ctx.font = "12px sans-serif";
    ctx.fillStyle = "#333333";
    const tagStr = truncate(opts.tags.join(", "), 25);
    ctx.fillText(tagStr, infoX, qrY + 8);
  }

  // Date
  if (opts.date) {
    ctx.font = "11px sans-serif";
    ctx.fillStyle = "#666666";
    ctx.fillText(opts.date, infoX, qrY + 28);
  }

  // Note ID / short URL
  ctx.font = "9px sans-serif";
  ctx.fillStyle = "#888888";
  const urlText = truncate(opts.url || "", 45);
  ctx.fillText(urlText, padding, height - padding - 4);

  return canvas;
}

function loadImage(url) {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => resolve(img);
    img.onerror = () => reject(new Error("Failed to load image: " + url));
    img.src = url;
  });
}

function truncate(str, maxLen) {
  return str.length > maxLen ? str.substring(0, maxLen - 1) + "\u2026" : str;
}

// ============================================================================
// PRINT MANAGER (ties everything together)
// ============================================================================

/**
 * High-level print manager used by the UI.
 *
 * Manages printer connection state and provides a simple API
 * for the notes view template to call.
 */
class LabelPrintManager {
  constructor() {
    this.printer = new NiimbotPrinter();
    this._statusEl = null;
    this._previewEl = null;
  }

  /**
   * Initialize with DOM element IDs for status and preview.
   */
  init(statusElId, previewElId) {
    this._statusEl = document.getElementById(statusElId);
    this._previewEl = document.getElementById(previewElId);
    this.printer.onStatus((msg) => {
      if (this._statusEl) this._statusEl.textContent = msg;
    });
  }

  /**
   * Connect to printer (opens browser BLE picker).
   */
  async connect() {
    await this.printer.connect();
    return this.printer.connected;
  }

  /**
   * Disconnect from printer.
   */
  disconnect() {
    this.printer.disconnect();
  }

  get isConnected() {
    return this.printer.connected;
  }

  /**
   * Render label preview and optionally print it.
   *
   * @param {Object} noteData - { qrUrl, title, url, date, tags }
   * @param {Object} printOpts - { density, labelType, quantity, labelWidth, labelHeight }
   * @param {boolean} doPrint - If true, send to printer after rendering
   */
  async renderAndPrint(noteData, printOpts = {}, doPrint = false) {
    const canvas = await renderLabel({
      qrImageUrl: noteData.qrUrl,
      title: noteData.title,
      url: noteData.url,
      date: noteData.date,
      tags: noteData.tags,
      width: printOpts.labelWidth || 320,
      height: printOpts.labelHeight || 240,
    });

    // Show preview
    if (this._previewEl) {
      this._previewEl.innerHTML = "";
      canvas.style.maxWidth = "100%";
      canvas.style.border = "1px solid #ddd";
      canvas.style.borderRadius = "4px";
      this._previewEl.appendChild(canvas);
    }

    if (doPrint) {
      if (!this.printer.connected) {
        throw new Error("Connect to printer first");
      }
      await this.printer.printLabel(canvas, {
        density: printOpts.density || 3,
        labelType: printOpts.labelType || 1,
        quantity: printOpts.quantity || 1,
      });
    }

    return canvas;
  }
}

// Export for use in template
window.LabelPrintManager = LabelPrintManager;
