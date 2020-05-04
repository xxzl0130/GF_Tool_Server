function FindProxyForURL(url, host)
{
    if ((shExpMatch(host, "gf-*-cn-zs-game-0001.ppgame.com") || shExpMatch(host, "s*.gw.gf.ppgame.com") 
    || shExpMatch(host, "s*.ios.gf.ppgame.com") || shExpMatch(host, "sn-game.txwy.tw") 
    || shExpMatch(host, "*girlfrontline.co.kr") || shExpMatch(host, "gfcn-game.*.sunborngame.com")))
        return "PROXY ar.xuanxuan.tech:8081";
    return "DIRECT";
}

