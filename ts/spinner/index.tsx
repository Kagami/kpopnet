/**
 * Simple pure CSS spinner.
 *
 * Based on https://loading.io/css/ example.
 *
 * @module kpopnet/spinner
 */

import * as cx from "classnames";
import { h } from "preact";
import "./index.less";

interface SpinnerProps {
  center?: boolean;
  large?: boolean;
}

function Spinner({ center, large }: SpinnerProps = {}) {
  return (
    <div class={cx("spinner", center && "spinner_centered", large && "spinner_2x")}>
      {Array(12).fill(
        <div class="spinner__blade" />,
      )}
    </div>
  );
}

export default Spinner;
