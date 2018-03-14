/**
 * Filterable idol list.
 */

import { Component, h } from "preact";
import "./index.css";

class Idol extends Component<any, any> {
  private url = "";
  public render() {
    return (
      <div class="idol">
        <img
          class="idol__img"
          src={this.url}
          draggable={0 as any}
          onDragStart={this.handleDragStart}
        />
        <div class="idol__info">
          <p class="idol__info-line">Stage name: Eunwoo</p>
          <p class="idol__info-line">Real name: Jung Eunwoo (정은우)</p>
          <p class="idol__info-line">Position: Main Vocalist</p>
          <p class="idol__info-line">Birthday: July 1, 1998</p>
          <p class="idol__info-line">Zodiac sign: Cancer</p>
          <p class="idol__info-line">Height: 166.6 cm</p>
          <p class="idol__info-line">Weight: 48 kg</p>
          <p class="idol__info-line">Blood Type: B</p>
        </div>
      </div>
    );
  }
  private handleDragStart = (e: DragEvent) => {
    e.preventDefault();
  }
}

class IdolList extends Component<any, any> {
  public render() {
    return <Idol/>;
  }
}

export default IdolList;
