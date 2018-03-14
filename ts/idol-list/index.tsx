/**
 * Filterable idol list.
 */

import { Component, h } from "preact";
import { getIdolPreviewUrl, Idol, Profiles } from "../api";
import "./index.css";

interface ItemProps {
  info: Idol;
}

class IdolItem extends Component<ItemProps, any> {
  public render({ info }: ItemProps) {
    return (
      <section class="idol">
        <img
          class="idol__preview"
          src={getIdolPreviewUrl(info.id)}
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
      </section>
    );
  }
  private handleDragStart = (e: DragEvent) => {
    e.preventDefault();
  }
}

interface ListProps {
  profiles: Profiles;
  query: string;
}

class IdolList extends Component<ListProps, any> {
  public render({ profiles }: ListProps) {
    const idols = profiles.idols.slice(0, 10);
    return (
      <article class="idol-list">
        {idols.map((info: Idol) =>
          <IdolItem
            key={info.id}
            info={info}
          />,
        )}
      </article>
    );
  }
}

export default IdolList;
