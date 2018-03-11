import { Component, h, render } from "preact";
import "./index.css";

class Index extends Component<any, any> {
  public render() {
    return <h1>test</h1>;
  }
}

render(<Index/>, document.querySelector(".app"));
