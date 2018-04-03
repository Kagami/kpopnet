import { Component, h } from "preact";
import "./index.less";

const ALLOWED_MIMES = new Set(["image/jpeg", "image/png"]);
const MAX_FILE_SIZE = 5 * 1024 * 1024;
const MIN_DIMENSION = 300;
const MAX_DIMENTION = 5000;

function validateImage(file: File): Promise<File> {
  return new Promise((resolve, reject) => {
    if (!ALLOWED_MIMES.has(file.type)) {
      throw new Error("Only JPEG and PNG allowed");
    }
    if (file.size > MAX_FILE_SIZE) {
      throw new Error("Max file size is 5MB");
    }
    const src = URL.createObjectURL(file);
    const img = new Image();
    img.onload = () => {
      const { width, height } = img;
      if (Math.min(width, height) < MIN_DIMENSION) {
        reject(new Error("Minimal resolution is 300x300"));
        return;
      }
      if (Math.max(width, height) > MAX_DIMENTION) {
        reject(new Error("Minimal resolution is 5000x5000"));
        return;
      }
      resolve(file);
    };
    img.onerror = () => {
      reject(new Error("Cannot load image"));
    };
    img.src = src;
  });
}

class Dropzone extends Component<any, any> {
  private fileEl: HTMLInputElement = null;
  public render() {
    return (
      <div
        class="dropzone"
        onClick={this.handleClick}
        onDragOver={this.handleDragOver}
        onDrop={this.handleDrop}
      >
        Click/drop photo of idol here
        <input
          ref={(f) => this.fileEl = f as HTMLInputElement}
          type="file"
          accept="image/*"
          class="dropzone__file"
          onChange={this.handleInputChange}
        />
      </div>
    );
  }
  private handleClick = () => {
    this.fileEl.click();
  }
  private handleInputChange = () => {
    const files = this.fileEl.files;
    if (files.length) {
      this.handleFile(files[0]);
    }
    this.fileEl.value = "";  // Allow to select same file again
  }
  private handleDragOver = (e: DragEvent) => {
    e.preventDefault();
  }
  private handleDrop = (e: DragEvent) => {
    e.preventDefault();
    const files = e.dataTransfer.files;
    if (files.length) {
      this.handleFile(files[0]);
    }
  }
  private handleFile(file: File) {
    validateImage(file).then(this.props.onChange, (err) => {
      alert(err);
    });
  }
}

export default Dropzone;
