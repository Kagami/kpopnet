/**
 * Application entry point.
 */

// tslint:disable-next-line:no-reference
/// <reference path="./index.d.ts" />

import { Component, h, render } from "preact";
import { BandMap, getBandMap, getIdolMap, getProfiles, IdolMap, Profiles } from "../api";
import Dropzone from "../dropzone";
import IdolList from "../idol-list";
import Recognizer from "../recognizer";
import Search from "../search";
import "./index.less";

interface IndexState {
  loading: boolean;
  query: string;
  file?: File;
}

class Index extends Component<{}, IndexState> {
  private profiles: Profiles = null;
  private bandMap: BandMap = null;
  private idolMap: IdolMap = null;
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
      this.idolMap = getIdolMap(profiles);
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
            query={query}
            loading={loading}
            disabled={!!file}
            onChange={this.handleSearch}
          />
          {!loading && !file && query &&
            <IdolList
              profiles={this.profiles}
              bandMap={this.bandMap}
              query={query}
            />
          }
          {!file && !query &&
            <Dropzone
              disabled={loading}
              onChange={this.handleFile}
            />
          }
          {file &&
            <Recognizer
              file={file}
              onFound={this.handleFound}
            />
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
  private handleFound = (idolId: string) => {
    // Everything must exist unless in a very rare case (e.g. new idols
    // was added after page load and user uploaded image with them.)
    const idol = this.idolMap.get(idolId);
    const iname = idol.name;
    const bname = this.bandMap.get(idol.band_id).name;
    const query = `${iname} band:${bname}`;
    this.setState({query, file: null});
  }
}

render(<Index/>, document.body);
