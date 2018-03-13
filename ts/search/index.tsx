import { Component, h } from "preact";
import "./index.css";

// https://loading.io/css/
function Spinnder() {
  return (
    <div class="spinner">
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
      <div class="spinner__lobe" />
    </div>
  );
}

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
          autofocus
        />
        {loading && <Spinnder/>}
      </div>
    );
  }
  private focus() {
    if (this.inputEl) {
      this.inputEl.focus();
    }
  }
}

export default Search;
