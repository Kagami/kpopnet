import { Component, h, render } from "preact";
import { getProfiles, Profiles } from "../api";
import IdolList from "../idol-list";
import Search from "../search";
import "./index.css";

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

class Index extends Component<any, any> {
  private profiles: Profiles = null;
  constructor() {
    super();
    this.state = {
      loading: true,
      search: "",
      file: null,
    };
  }
  public componentDidMount() {
    // FIXME(Kagami): Error handling.
    getProfiles().then((profiles) => {
      this.profiles = profiles;
      this.setState({loading: false});
    });
  }
  public render({}, { loading, search, file }: any) {
    return (
      <main class="index">
        <div class="index__inner">
          <Search
            loading={loading}
            onChange={this.handleSearch}
          />
          {!file && <Dropzone onLoad={this.handleLoad} />}
          {search && <IdolList profiles={this.profiles} search={search} />}
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
      </main>
    );
  }
  private handleLoad = (file: File) => {
    this.setState({file});
  }
  private handleSearch = (search: string) => {
    this.setState({search});
  }
}

render(<Index/>, document.body);
