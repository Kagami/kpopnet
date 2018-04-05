import * as cx from "classnames";
import { Component, h } from "preact";
import { hook, HOOKS, trigger, unhook } from "../hooks";
import "./index.less";

interface Alert {
  id?: number;
  title?: string;
  message: string;
  sticky?: boolean;
  closing?: boolean;
}

interface AlertsState {
  alerts: Alert[];
}

class Alerts extends Component<{}, AlertsState> {
  public state: AlertsState = {
    alerts: [],
  };
  private id = 0;
  public componentDidMount() {
    hook(HOOKS.showAlert, this.show);
  }
  public componentWillUnmount() {
    unhook(HOOKS.showAlert, this.show);
  }
  public render({}, { alerts }: AlertsState) {
    return (
      <aside class="alerts">
        {alerts.map(this.renderAlert)}
      </aside>
    );
  }
  private show = (a: Alert) => {
    a = Object.assign({}, a, {id: this.id++, closing: false});
    const alerts = [a].concat(this.state.alerts);
    this.setState({alerts});
    if (!a.sticky) {
      setTimeout(this.makeClose(a.id), 4000);
    }
  }
  private makeClose(id: number) {
    return () => {
      const alerts = this.state.alerts.map((a) =>
        a.id === id ? {...a, closing: true} : a,
      );
      this.setState({alerts});
      setTimeout(() => {
        // tslint:disable-next-line:no-shadowed-variable
        const alerts = this.state.alerts.filter((a) => a.id !== id);
        this.setState({alerts});
      }, 1000);
    };
  }
  private renderAlert = ({ id, title, message, closing }: Alert) => {
    return (
      <article class={cx("alert", closing && "alert_closing")} key={id.toString()}>
        <a class="alert-close-control" onClick={this.makeClose(id)}>✖</a>
        {this.renderTitle(title)}
        <section class="alert-message">{message}</section>
      </article>
    );
  }
  private renderTitle(title: string) {
    return title ? <header class="alert-title">{title}</header> : null;
  }
}

export default Alerts;

export function showAlert(a: Alert | Error | string | [string, string]) {
  if (typeof a === "string") {
    a = {message: a};
  } else if (a instanceof Error) {
    a = {message: a.message};
  } else if (Array.isArray(a)) {
    a = {title: a[0], message: a[1]};
  }
  trigger(HOOKS.showAlert, a);
}
