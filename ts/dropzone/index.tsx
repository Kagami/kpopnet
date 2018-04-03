import { Component, h } from "preact";
import "./index.less";

const ALLOWED_MIMES = new Set(["image/jpeg", "image/png"]);
const MAX_FILE_SIZE = 5 * 1024 * 1024;

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
    if (!ALLOWED_MIMES.has(file.type)) {
      alert("Only JPEG and PNG allowed");
      return;
    }
    if (file.size > MAX_FILE_SIZE) {
      alert("Max file size is 5MB");
      return;
    }
    this.props.onChange(file);
  }
}

export default Dropzone;
