# Simple VUE project to switch to SBIT network via Metamask

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

### wallet_addEthereumChain
```
// request account access
window.sbit.request({ method: 'eth_requestAccounts' })
    .then(() => {
        // add chain
        window.sbit.request({
            method: "wallet_addEthereumChain",
            params: [{
                {
                    chainId: '0x22B9',
                    chainName: 'Sbit Testnet',
                    rpcUrls: ['https://localhost:22402'],
                    blockExplorerUrls: ['https://testnet.sbit.dev/'],
                    iconUrls: [
                        'https://mainnet.sbit.dev/images/metamask_icon.svg',
                        'https://mainnet.sbit.dev/images/metamask_icon.png',
                    ],
                    nativeCurrency: {
                        decimals: 18,
                        symbol: 'SBIT',
                    },
                }
            }],
        }
    });
```

# Known issues
- Metamask requires https for `rpcUrls` so that must be enabled
  - Either directly through Janus with `--https-key ./path --https-cert ./path2` see [SSL](../README.md#ssl)
  - Through the Makefile `make docker-configure-https && make run-janus-https`
  - Or do it yourself with a proxy (eg, nginx)
