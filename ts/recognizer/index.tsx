import { Component, h } from "preact";
import { recognizeIdol } from "../api";
import Spinner from "../spinner";
import "./index.less";

interface RecognizerProps {
  file: File;
  onMatch: (idolId: string) => void;
  onError: (err: Error) => void;
}

class Recognizer extends Component<RecognizerProps, {}> {
  private imageUrl = "";
  constructor(props: RecognizerProps) {
    super(props);
    this.imageUrl = URL.createObjectURL(props.file);
  }
  public componentWillMount() {
    recognizeIdol(this.props.file, {prefix: API_PREFIX}).then(({ id }) => {
      this.props.onMatch(id);
    }, this.props.onError);
  }
  public render({ file }: RecognizerProps) {
    return (
      <div class="recognizer recognizer_loading">
        <img
          class="recognizer__preview"
          src={this.imageUrl}
          draggable={0 as any}
          onDragStart={this.handleDrag}
        />
        <Spinner center large />
      </div>
    );
  }
  private handleDrag = (e: DragEvent) => {
    e.preventDefault();
  }
}

export default Recognizer;
