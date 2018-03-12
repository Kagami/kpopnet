import { Component, h, render } from "preact";
import "./index.css";

class Search extends Component<any, any> {
  public render() {
    return (
      <input
        class="search"
        placeholder="Search for idol or band"
        autofocus
      />
    );
  }
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
    this.props.onLoad(file);
  }
}

class Idol extends Component<any, any> {
  private url: string = null;
  public componentWillMount() {
    this.url = URL.createObjectURL(this.props.file);
  }
  public componentWillUnmount() {
    URL.revokeObjectURL(this.url);
  }
  public render() {
    return (
      <div class="idol">
        <img
          class="idol__img"
          src={this.url}
          draggable={0 as any}
          onDragStart={this.handleDragStart}
        />
        <div class="idol__info">
          <p class="idol__info-line">Stage name: Eunwoo</p>
          <p class="idol__info-line">Real name: Jung Eunwoo (정은우)</p>
          <p class="idol__info-line">Position: Main Vocalist</p>
          <p class="idol__info-line">Birthday: July 1, 1998</p>
          <p class="idol__info-line">Zodiac sign: Cancer</p>
          <p class="idol__info-line">Height: 166.6 cm</p>
          <p class="idol__info-line">Weight: 48 kg</p>
          <p class="idol__info-line">Blood Type: B</p>
        </div>
      </div>
    );
  }
  private handleDragStart = (e: DragEvent) => {
    e.preventDefault();
  }
}

class Index extends Component<any, any> {
  constructor() {
    super();
    this.state = {
      file: null,
    };
  }
  public render({}, { file }: any) {
    return (
      <div class="index">
        <div class="index__inner">
          <Search />
          {!file && <Dropzone onLoad={this.handleLoad} />}
          {file && <Idol file={file} />}
        </div>
        <footer class="footer">
          <div class="footer__inner">
            <a class="footer__link" href="https://kpop.re/">
              Kpop.re
            </a>
            <a class="footer__link" href="https://github.com/Kagami/kpopnet">
              Source code
            </a>
            <a class="footer__link" href="https://github.com/Kagami/kpopnet/issues">
              Feedback
            </a>
          </div>
        </footer>
      </div>
    );
  }
  private handleLoad = (file: File) => {
    this.setState({file});
  }
}

render(<Index/>, document.body);
