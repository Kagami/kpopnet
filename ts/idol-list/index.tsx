/**
 * Filterable idol list.
 */

import { Component, h } from "preact";
import {
  Band, BandMap, getIdolPreviewUrl, Idol, Profiles, renderIdol,
} from "../api";
import "./index.less";

interface ItemProps {
  idol: Idol;
  band: Band;
}

class IdolItem extends Component<ItemProps, any> {
  public shouldComponentUpdate() {
    return false;
  }
  public render({ idol, band }: ItemProps) {
    return (
      <section class="idol">
        <img
          class="idol__preview"
          src={getIdolPreviewUrl(idol.id)}
          draggable={0 as any}
          onDragStart={this.handleDragStart}
        />
        <div class="idol__info">
          {renderIdol(idol, band).map(([key, val]) =>
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
  bandMap: BandMap;
  query: string;
}

class IdolList extends Component<ListProps, any> {
  public shouldComponentUpdate(nextProps: ListProps) {
    return this.props.query !== nextProps.query;
  }
  public render({ profiles, bandMap }: ListProps) {
    const idols = profiles.idols.slice(0, 10);
    return (
      <article class="idol-list">
        {idols.map((idol) =>
          <IdolItem
            key={idol.id}
            idol={idol}
            band={bandMap.get(idol.band_id).band}
          />,
        )}
      </article>
    );
  }
}

export default IdolList;
