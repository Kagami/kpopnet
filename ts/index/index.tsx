import { Component, h, render } from "preact";
import { BandMap, getBandMap, getProfiles, Profiles } from "../api";
import Dropzone from "../dropzone";
import IdolList from "../idol-list";
import Recognizer from "../recognizer";
import Search from "../search";
import "./index.less";

declare const API_PREFIX: string;

class Index extends Component<any, any> {
  private profiles: Profiles = null;
  private bandMap: BandMap = null;
  constructor() {
    super();
    this.state = {
      loading: true,
      query: "",
      file: null,
    };
  }
  public componentDidMount() {
    getProfiles({prefix: API_PREFIX}).then((profiles) => {
      this.profiles = profiles;
      this.bandMap = getBandMap(profiles);
      this.setState({loading: false});
    }, (err) => {
      this.setState({loading: false});
      // TODO(Kagami): Something better.
      alert("Error getting profiles");
    });
  }
  public render({}, { loading, query, file }: any) {
    return (
      <main class="index">
        <div class="index__inner">
          <Search
            loading={loading}
            disabled={!!file}
            onChange={this.handleSearch}
          />
          {!file && !loading && query &&
            <IdolList
              profiles={this.profiles}
              bandMap={this.bandMap}
              query={query}
            />
          }
          {!file && !query &&
            <Dropzone onChange={this.handleFile} />
          }
          {file &&
            <Recognizer file={file} />
          }
        </div>
        <footer class="footer">
          <div class="footer__inner">
            <a class="footer__link" target="_blank" href="https://kpop.re/">
              Kpop.re
            </a>
            <a class="footer__link" target="_blank" href="https://github.com/Kagami/kpopnet">
              Source code
            </a>
            <a class="footer__link" target="_blank" href="https://github.com/Kagami/kpopnet/issues">
              Feedback
            </a>
          </div>
        </footer>
      </main>
    );
  }
  private handleFile = (file: File) => {
    this.setState({file});
  }
  private handleSearch = (query: string) => {
    this.setState({query});
  }
}

render(<Index/>, document.body);
