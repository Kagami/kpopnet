import { Component, h, render } from "preact";
import { BandMap, getBandMap, getProfiles, Profiles } from "../api";
import Dropzone from "../dropzone";
import IdolList from "../idol-list";
import Search from "../search";
import "./index.less";

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
    // FIXME(Kagami): Error handling.
    getProfiles().then((profiles) => {
      this.profiles = profiles;
      this.bandMap = getBandMap(profiles);
      this.setState({loading: false});
    });
  }
  public render({}, { loading, query, file }: any) {
    return (
      <main class="index">
        <div class="index__inner">
          <Search
            loading={loading}
            onChange={this.handleSearch}
          />
          {(!file && !query) &&
            <Dropzone onChange={this.handleFile} />
          }
          {(!loading && query) &&
            <IdolList
              profiles={this.profiles}
              bandMap={this.bandMap}
              query={query}
            />
          }
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
  private handleFile = (file: File) => {
    this.setState({file});
  }
  private handleSearch = (query: string) => {
    this.setState({query});
  }
}

render(<Index/>, document.body);
