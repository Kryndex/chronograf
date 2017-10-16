import React, {Component, PropTypes} from 'react'

const dummyLogs = [
  {
    ts: '2017-10-16T15:36:31.329-04:00',
    lvl: 'info',
    msg: 'created log session',
    service: 'sessions',
    id: 'aa87c0b6-26d0-484a-8f15-7445c2d5f386',
    'content-type': 'application/json',
    tags: 'nil',
  },
  {
    ts: '2017-10-16T15:36:31.329-04:00',
    lvl: 'error',
    msg: '2017/10/16 15:36:31 http: multiple response.WriteHeader calls\n',
    service: 'httpd_server_errors',
  },
  {
    ts: '2017-10-16T15:36:32.090-04:00',
    lvl: 'debug',
    msg: 'starting next batch query',
    service: 'kapacitor',
    task_master: 'main',
    task: 'batch',
    node: 'query1',
    query:
      "SELECT mean(usage_user) FROM telegraf.autogen.cpu WHERE time >= '2017-10-16T19:35:32.089928906Z' AND time < '2017-10-16T19:36:32.089928906Z'",
  },
  {
    ts: '2017-10-16T15:36:32.241-04:00',
    lvl: 'info',
    msg: 'http request',
    service: 'http',
    host: '::1',
    username: '-',
    start: '2017-10-16T15:36:32.241436147-04:00',
    method: 'POST',
    uri: '/write?consistency=&db=mydb&precision=ns&rp=autogen',
    protocol: 'HTTP/1.1',
    status: 204,
    referer: '-',
    'user-agent': 'InfluxDBClient',
    'request-id': '502cb71e-b2a9-11e7-8186-000000000000',
    duration: '59.135µs',
  },
  {
    ts: '2017-10-16T15:36:35.313-04:00',
    lvl: 'info',
    msg: 'http request',
    service: 'http',
    host: '::1',
    username: '-',
    start: '2017-10-16T15:36:35.312244484-04:00',
    method: 'POST',
    uri: '/write?consistency=&db=telegraf&precision=ns&rp=autogen',
    protocol: 'HTTP/1.1',
    status: 204,
    referer: '-',
    'user-agent': 'InfluxDBClient',
    'request-id': '52014892-b2a9-11e7-8187-000000000000',
    duration: '1.468473ms',
  },
  {
    ts: '2017-10-16T15:36:35.313-04:00',
    lvl: 'info',
    msg: 'point',
    service: 'kapacitor',
    task_master: 'main',
    task: 'log',
    node: 'log2',
    prefix: '',
    name: 'system',
    db: 'telegraf',
    rp: 'autogen',
    group: 'cluster_id=michaels-example-cluster,host=Michaels-MBP-2.router.edm',
    dimension_0: 'cluster_id',
    dimension_1: 'host',
    tag: {
      cluster_id: 'michaels-example-cluster',
      host: 'Michaels-MBP-2.router.edm',
    },
    field: {
      load15: 1.71,
      n_users: 8,
      load5: 1.78,
      n_cpus: 8,
      load1: 1.74,
    },
    time: '2017-10-16T19:36:35Z',
  },
  {
    ts: '2017-10-16T15:36:35.313-04:00',
    lvl: 'info',
    msg: 'point',
    service: 'kapacitor',
    task_master: 'main',
    task: 'log',
    node: 'log2',
    prefix: '',
    name: 'system',
    db: 'telegraf',
    rp: 'autogen',
    group: 'cluster_id=michaels-example-cluster,host=Michaels-MBP-2.router.edm',
    dimension_0: 'cluster_id',
    dimension_1: 'host',
    tag: {
      cluster_id: 'michaels-example-cluster',
      host: 'Michaels-MBP-2.router.edm',
    },
    field: {
      uptime_format: '11 days,  2:08',
      uptime: 958109,
    },
    time: '2017-10-16T19:36:35Z',
  },
  {
    ts: '2017-10-16T15:36:36.664-04:00',
    lvl: 'info',
    msg: 'http request',
    service: 'http',
    host: '::1',
    username: '-',
    start: '2017-10-16T15:36:36.663967463-04:00',
    method: 'POST',
    uri: '/write?consistency=&db=_internal&precision=ns&rp=monitor',
    protocol: 'HTTP/1.1',
    status: 204,
    referer: '-',
    'user-agent': 'InfluxDBClient',
    'request-id': '52cf8a43-b2a9-11e7-8188-000000000000',
    duration: '624.155µs',
  },
  {
    ts: '2017-10-16T15:36:40.313-04:00',
    lvl: 'info',
    msg: 'http request',
    service: 'http',
    host: '::1',
    username: '-',
    start: '2017-10-16T15:36:40.312324709-04:00',
    method: 'POST',
    uri: '/write?consistency=&db=telegraf&precision=ns&rp=autogen',
    protocol: 'HTTP/1.1',
    status: 204,
    referer: '-',
    'user-agent': 'InfluxDBClient',
    'request-id': '54fc3c33-b2a9-11e7-8189-000000000000',
    duration: '1.0192ms',
  },
  {
    ts: '2017-10-16T15:36:40.313-04:00',
    lvl: 'info',
    msg: 'point',
    service: 'kapacitor',
    task_master: 'main',
    task: 'log',
    node: 'log2',
    prefix: '',
    name: 'system',
    db: 'telegraf',
    rp: 'autogen',
    group: 'cluster_id=michaels-example-cluster,host=Michaels-MBP-2.router.edm',
    dimension_0: 'cluster_id',
    dimension_1: 'host',
    tag: {
      host: 'Michaels-MBP-2.router.edm',
      cluster_id: 'michaels-example-cluster',
    },
    field: {
      load15: 1.71,
      n_cpus: 8,
      load1: 1.68,
      load5: 1.77,
      n_users: 8,
    },
    time: '2017-10-16T19:36:40Z',
  },
]
class LogsTable extends Component {
  constructor(props) {
    super(props)
  }

  renderTable = () => {
    return (
      <table className="table logs-table">
        <thead>
          <tr>
            <th>Blargh</th>
            <th>Swoggle</th>
            <th>Horgles</th>
            <th>Chortle</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>sdfsdfsf</td>
            <td>sdfsdfsdf</td>
            <td>sdfsdfsdfjnsdfkjlnjsdjkfsf</td>
            <td>sdbnds</td>
          </tr>
          <tr>
            <td>sdfsdfsf</td>
            <td>sdfsdfsdf</td>
            <td>sdfsdfsdfjnsdfkjlnjsdjkfsf</td>
            <td>sdbnds</td>
          </tr>
        </tbody>
      </table>
    )
  }
  render() {
    const {isWidget} = this.props

    const output = isWidget
      ? this.renderTable()
      : <div className="logs-table-container">
          <div className="panel panel-minimal">
            <div className="panel-heading u-flex u-ai-center u-jc-space-between">
              <h2 className="panel-title">Logs</h2>
              <div className="filterthing">FILTER</div>
            </div>
            <div className="panel-body">
              {this.renderTable()}
            </div>
          </div>
        </div>

    return output
  }
}

const {bool} = PropTypes

LogsTable.propTypes = {
  isWidget: bool,
}

export default LogsTable