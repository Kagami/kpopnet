import * as cx from "classnames";
import { Component, h } from "preact";
import { recognizeIdol } from "../api";
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
  public componentWillMount() {
    recognizeIdol(this.props.file, {prefix: API_PREFIX}).then(({ id }) => {
      this.setState({loading: false});
      alert(`Recognized: ${id}`);
    }, (err) => {
      this.setState({loading: false});
      alert(`Error recognizing: ${err.message}`);
    });
  }
  public render({ file }: RecognizerProps, { loading }: RecognizerState) {
    return (
      <div class={cx("recognizer", loading && "recognizer_loading")}>
        <img
          class="recognizer__preview"
          src={this.imageUrl}
          draggable={0 as any}
          onDragStart={this.handleDrag}
        />
        {loading && <Spinner center large />}
      </div>
    );
  }
  private handleDrag = (e: DragEvent) => {
    e.preventDefault();
  }
}

export default Recognizer;
