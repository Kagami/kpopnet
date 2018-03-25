/**
 * Filterable idol list.
 */

// tslint:disable-next-line:no-reference
/// <reference path="./index.d.ts" />

import { Component, h } from "preact";
import {
  BandMap, getIdolPreviewUrl, Idol,
  Profiles, renderIdol, searchIdols,
} from "../api";
import "./index.less";
import previewFallbackUrl from "./no-preview.svg";

declare const FILE_PREFIX: string;

interface ItemProps {
  idol: Idol;
  bandMap: BandMap;
}

class IdolItem extends Component<ItemProps, any> {
  public shouldComponentUpdate() {
    return false;
  }
  public render({ idol, bandMap }: ItemProps) {
    const opts = {prefix: FILE_PREFIX, fallback: previewFallbackUrl};
    const previewUrl = getIdolPreviewUrl(idol, opts);
    const style = {backgroundImage: `url(${previewUrl})`};
    return (
      <section class="idol">
        <div
          class="idol__preview"
          style={style}
        />
        <div class="idol__info">
          {renderIdol(idol, bandMap).map(([key, val]) =>
            <p class="idol__info-line">
              <span class="idol__info-key">{key}</span>
              <span class="idol__info-val">{val}</span>
            </p>,
          )}
        </div>
      </section>
    );
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
  public render({ query, profiles, bandMap }: ListProps) {
    const idols = searchIdols(query, profiles, bandMap).slice(0, 20);
    if (!idols.length) return this.renderEmpty();
    return (
      <article class="idols">
        {idols.map((idol) =>
          <IdolItem
            key={idol.id}
            idol={idol}
            bandMap={bandMap}
          />,
        )}
      </article>
    );
  }
  public renderEmpty() {
    return (
      <article class="idols idols_empty">
        No results
      </article>
    );
  }
}

export default IdolList;
