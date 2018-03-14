import { Component, h } from "preact";
import Spinner from "../spinner";
import "./index.css";

class Search extends Component<any, any> {
  private inputEl: HTMLInputElement = null;
  public componentDidUpdate(prevProps: any) {
    if (prevProps.loading && !this.props.loading) {
      this.focus();
    }
  }
  public render({ loading }: any) {
    return (
      <div class="search">
        <input
          ref={(i) => this.inputEl = i as HTMLInputElement}
          class="search__input"
          placeholder="Search for idol or band"
          disabled={loading}
          onInput={this.handleChange}
        />
        {loading && <Spinner/>}
      </div>
    );
  }
  private focus() {
    if (this.inputEl) {
      this.inputEl.focus();
    }
  }
  private handleChange = (e: Event) => {
    this.props.onChange((e.target as HTMLInputElement).value);
  }
}

export default Search;
