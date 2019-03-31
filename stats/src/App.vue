<template>
    <div id="app">
        <b-card>
            <b-row>
                <b-col>
                    <h3>
                        <b-button @click="showAlert" variant="light" class="m-1">
                            <img src="./assets/refresh-ccw.svg" width="16" :hidden="refreshing">
                            <img src="./assets/refresh-ccw.svg" width="16" :hidden="!refreshing">
                        </b-button>
                        RPC-API Requests
                    </h3>
                </b-col>
                <b-col cols="3">
                    <b-alert
                            :show="dismissCountDown"
                            dismissible
                            variant="primary"
                            @dismissed="dismissCountDown=0"
                            @dismiss-count-down="countDownChanged"
                    ><div align="center">Data Refreshed</div>
                    </b-alert>
                </b-col>
            </b-row>
            <b-card-body>
                <b-tabs content-class="mt-3" pills card >
                    <b-tab title="Network" active>
                        <span class="my-auto">
                            <div height="800px" width="100%" align="center" class="bg-gradient-dark" id="attackNetwork">
                                <h3 align="center" class="large text-white">Graph Loading ...</h3>
                                <p align="center" class="text-white">[ ether Transfer Attempts by Destination Ethereum Account &#8227; IP Address ]</p>
                                <p align="center" class="small text-white">click to highlight connections, double click to get details</p>
                            </div>
                        </span>
                    </b-tab>
                    <b-tab title="IP Address">
                        <b-table striped hover :items="summary"/>
                        <table class="table table-striped">
                            <thead class="thead-dark">
                            <tr>
                                <th>IP</th>
                                <th>Methods</th>
                                <th>Addresses (To:)</th>
                                <th>Last Seen</th>
                            </tr>
                            </thead>
                            <tr v-for="(info, ip) in stats.ips">
                                <td class="badge-light text-lg-center my-auto">
                                    <b-link :href="shodanLink(ip)" target="_new" class="my-auto text-dark">{{ip}}</b-link>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, method) in info.methods">
                                            <td><b-link :href="docsLink(method)" target="_new" class="my-auto text-dark">{{method}}</b-link></td>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, address) in info.accounts">
                                            <td><b-link :href="etherscanLink(address)" target="_new" class="my-auto text-dark">{{address}}</b-link></td>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    {{info.last_seen}}
                                </td>
                            </tr>
                        </table>
                    </b-tab>
                    <b-tab title="Method">
                        <b-table striped hover :items="summary"/>
                        <table class="table table-striped">
                            <thead class="thead-dark">
                            <tr>
                                <th>Method</th>
                                <th>IPs</th>
                                <th>Addresses (To:)</th>
                                <th>Last Seen</th>
                            </tr>
                            </thead>
                            <tr v-for="(info, method) in stats.methods">
                                <td class="badge-light text-lg-center my-auto">
                                    <b-link :href="docsLink(method)" target="_new" class="my-auto text-dark">{{method}}</b-link>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, ip) in info.ips">
                                            <td>
                                                <b-link :href="shodanLink(ip)" target="_new" class="my-auto text-dark">{{ip}}</b-link>
                                            </td>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, address) in info.accounts">
                                            <td><b-link :href="etherscanLink(address)" target="_new" class="my-auto text-dark">{{address}}</b-link></td>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    {{info.last_seen}}
                                </td>
                            </tr>
                        </table>
                    </b-tab>
                    <b-tab title="Account">
                        <b-table striped hover :items="summary"/>
                        <table class="table table-striped">
                            <thead class="thead-dark">
                            <tr>
                                <th>Account (To:)</th>
                                <th>IPs</th>
                                <th>Methods</th>
                                <th>Last Seen</th>
                            </tr>
                            </thead>
                            <tr v-for="(info, address) in stats.Addresses">
                                <td class="badge-light text-lg-center my-auto">
                                    <b-link :href="etherscanLink(address)" target="_new" class="my-auto text-dark">{{address}}</b-link>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, ip) in info.ips">
                                            <td>
                                                <b-link :href="shodanLink(ip)" target="_new" class="my-auto text-dark">{{ip}}</b-link>
                                            </td>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    <table>
                                        <tr v-for="(times, method) in info.methods">
                                            <td><b-link :href="docsLink(method)" target="_new" class="my-auto text-dark">{{method}}</b-link>
                                            <td>{{times}}</td>
                                        </tr>
                                    </table>
                                </td>
                                <td>
                                    {{info.last_seen}}
                                </td>
                            </tr>
                        </table>
                    </b-tab>
                    <b-tab title="Raw Data">
                        <p align="center">
                            <a href="/data/attacks.json">Download Full attack log</a> | <a href="/data/ip-account.gv">Download Graph DOT File</a> | <a href="/data/addrs.json">Download IP Address Data</a> | <a href="/data/stats.json">Download Summary Data</a><br />
                        </p>
                        <h5>Summary Data:</h5>
                        <div class="badge-light">
                            <pre><code>{{stats}}</code></pre>
                        </div>
                    </b-tab>
                </b-tabs>
            </b-card-body>
        </b-card>
        <span></span>
    </div>
</template>

<script>
import axios from 'axios';
import vis from 'vis';
export default {
  name: 'app',
  data () {
    return {
      stats: null,
        summary: [],
        refreshing: false,
        dismissSecs: 3,
        dismissCountDown: 0,
        showDismissibleAlert: false,
        network: null
    }
  },
  mounted() {
     this.refresh()
  },
    methods: {
        refresh() {
            this.refreshing = true;
            axios
                .get('/data/stats.json')
                .then(response => (
                    this.stats = response.data,
                        this.summary = this.summarize()));
            this.fetchGraph();
            this.refreshing = false;
        },
        summarize() {
            return [{
                'Unique_IP_Addresses': Object.keys(this.stats.ips).length,
                'Unique_Destination_Addresses': Object.keys(this.stats.Addresses).length,
                'Unique_API_Methods_Called': Object.keys(this.stats.methods).length
            }]
        },
        countDownChanged(dismissCountDown) {
            this.dismissCountDown = dismissCountDown
        },
        showAlert() {
            this.refresh();
            this.dismissCountDown = this.dismissSecs
        },
        etherscanLink(address) {
            return 'https://etherscan.io/address/' + address
        },
        shodanLink(address) {
            return 'https://www.shodan.io/search?query=' + address
        },
        docsLink(apiMethod) {
            return 'https://github.com/ethereum/wiki/wiki/JSON-RPC#' + apiMethod
        },
        fetchGraph() {
            axios
                .get('/data/ip-account.gv')
                .then(response => (
                    this.buildGraph(response.data)));
        },
        buildGraph(dotGv) {
            let parsedData = vis.network.convertDot(dotGv);
            let data = {
                nodes: parsedData.nodes,
                edges: parsedData.edges
            };
            let options = {
                autoResize: false,
                height: '800px',
                width: '100%',
                clickToUse: false,
                layout:{
                    randomSeed: 9999
                },
                edges:{
                    color:{
                        highlight: '#ffffff',
                        hover: '#0088ff',
                    }
                },
                nodes: {
                    shadow: true,
                    scaling: {
                        label: {
                            min: 10,
                            max: 50
                        }
                    },
                    borderWidthSelected: 3,
                    borderWidth: 0.3,
                    margin: 1,
                    color:{
                        border: '#000000',
                        highlight:{
                            border: '#a33434',
                            background: '#00960c'
                        }
                    }
                },
                physics: {
                   stabilization:{
                       enabled: false
                   },
                   solver: 'barnesHut',
                   barnesHut:{
                       gravitationalConstant: -2000,
                       centralGravity: 3.5,
                       avoidOverlap: 0.75
                   },
                   timestep: 0.15,
                   maxVelocity: 1250,
                   minVelocity: 500
                },
                interaction: {
                    hover: true
                }
            };
            let container = document.getElementById('attackNetwork');
            this.network = new vis.Network(container, data, options);
            this.network.on("doubleClick", function (params) {
                if (params.nodes.length === 1) {
                    if (params.nodes[0].search(/^\d{1,3}\./)) {
                        window.open('https://etherscan.io/address/' + params.nodes[0], '_blank');
                    } else {
                        window.open('https://www.shodan.io/search?query=' + params.nodes[0], '_blank');
                    }
                }
            });
        }
    }
}
</script>

<style>
</style>
