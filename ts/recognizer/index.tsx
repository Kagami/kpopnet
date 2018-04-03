import * as cx from "classnames";
import { Component, h } from "preact";
import Spinner from "../spinner";
import "./index.less";

interface RecognizerProps {
  file: File;
}

interface RecognizerState {
  loading: boolean;
}

class Recognizer extends Component<RecognizerProps, any> {
  private imageUrl = "";
  constructor(props: RecognizerProps) {
    super(props);
    this.imageUrl = URL.createObjectURL(props.file);
    this.state = {
      loading: true,
    };
  }
  public render({ file }: RecognizerProps, { loading }: RecognizerState) {
    const style = {backgroundImage: `url(${this.imageUrl})`};
    return (
      <div class={cx("recognizer", loading && "recognizer_loading")}>
        <div
          class="recognizer__preview"
          style={style}
        />
        {loading && <Spinner center large />}
      </div>
    );
  }
}

export default Recognizer;
