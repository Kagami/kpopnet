import { Component, h } from "preact";
import Spinner from "../spinner";
import "./index.less";

interface SearchProps {
  query: string;
  loading: boolean;
  disabled: boolean;
  onChange: (query: string) => void;
}

class Search extends Component<SearchProps, {}> {
  private inputEl: HTMLInputElement = null;
  public componentDidUpdate({ loading }: SearchProps) {
    if (loading && !this.props.loading) {
      this.focus();
    }
  }
  public render({ query, loading, disabled }: SearchProps) {
    return (
      <div class="search">
        <input
          ref={(i) => this.inputEl = i as HTMLInputElement}
          class="search__input"
          value={query}
          maxLength={40}
          placeholder="Search for idol or band"
          disabled={loading || disabled}
          onInput={this.handleInputChange}
        />
        {this.renderClearButton()}
        {loading && <Spinner/>}
      </div>
    );
  }
  private renderClearButton() {
    if (!this.props.query) return null;
    return (
      <span class="search__clear-control" onClick={this.handleClearClick}>
        ✖
      </span>
    );
  }
  private focus = () => {
    if (this.inputEl) {
      this.inputEl.focus();
    }
  }
  private handleInputChange = () => {
    this.props.onChange(this.inputEl.value);
  }
  private handleClearClick = () => {
    this.props.onChange("");
    setTimeout(this.focus);
  }
}

export default Search;
