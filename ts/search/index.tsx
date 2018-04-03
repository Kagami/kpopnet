import { Component, h } from "preact";
import Spinner from "../spinner";
import "./index.less";

interface SearchProps {
  loading: boolean;
  disabled: boolean;
  onChange: (query: string) => void;
}

class Search extends Component<SearchProps, any> {
  private inputEl: HTMLInputElement = null;
  public componentDidUpdate(prevProps: any) {
    if (prevProps.loading && !this.props.loading) {
      this.focus();
    }
  }
  public render({ loading, disabled }: SearchProps) {
    return (
      <div class="search">
        <input
          ref={(i) => this.inputEl = i as HTMLInputElement}
          class="search__input"
          placeholder="Search for idol or band"
          disabled={loading || disabled}
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
  private handleChange = () => {
    this.props.onChange(this.inputEl.value);
  }
}

export default Search;
