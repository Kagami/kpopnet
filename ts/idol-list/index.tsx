/**
 * Filterable idol list.
 */

import { Component, h } from "preact";
import { getIdolPreviewUrl, Idol, Profiles, showIdol } from "../api";
import "./index.less";

interface ItemProps {
  idol: Idol;
}

class IdolItem extends Component<ItemProps, any> {
  public render({ idol }: ItemProps) {
    return (
      <section class="idol">
        <img
          class="idol__preview"
          src={getIdolPreviewUrl(idol.id)}
          draggable={0 as any}
          onDragStart={this.handleDragStart}
        />
        <div class="idol__info">
          {showIdol(idol).map(([key, val]) =>
            <p class="idol__info-line">{key}: {val}</p>,
          )}
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
        {idols.map((idol: Idol) =>
          <IdolItem
            key={idol.id}
            idol={idol}
          />,
        )}
      </article>
    );
  }
}

export default IdolList;
