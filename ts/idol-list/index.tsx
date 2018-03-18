/**
 * Filterable idol list.
 */

import { Component, h } from "preact";
import {
  Band, BandMap, getIdolPreviewUrl, Idol, Profiles,
  renderIdol, searchIdols,
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
    const previewUrl = getIdolPreviewUrl(idol.id);
    const style = {backgroundImage: `url(${previewUrl})`};
    return (
      <section class="idol">
        <div
          class="idol__preview"
          style={style}
        />
        <div class="idol__info">
          {renderIdol(idol, band).map(([key, val]) =>
            <p class="idol__info-line">
              <span class="idol__info-key">{key}:</span>
              <span class="idol__info-val">{val}</span>
            </p>,
          )}
        </div>
      </section>
    );
  }
}

function IncompleteList() {
  return (
    <div class="idols__incomplete">
      <div>Some results were skipped,</div>
      <div>please clarify request</div>
    </div>
  );
}

interface ListProps {
  profiles: Profiles;
  bandMap: BandMap;
  query: string;
}

class IdolList extends Component<ListProps, any> {
  private MAX_ITEMS_COUNT = 10;
  public shouldComponentUpdate(nextProps: ListProps) {
    return this.props.query !== nextProps.query;
  }
  public render({ query, profiles, bandMap }: ListProps) {
    let idols = searchIdols(query, profiles, bandMap);
    let complete = true;
    if (idols.length > this.MAX_ITEMS_COUNT) {
      idols = idols.slice(0, this.MAX_ITEMS_COUNT);
      complete = false;
    }
    return (
      <article class="idols">
        {idols.map((idol) =>
          <IdolItem
            key={idol.id}
            idol={idol}
            band={bandMap.get(idol.band_id).band}
          />,
        )}
        {!complete && <IncompleteList />}
      </article>
    );
  }
}

export default IdolList;
