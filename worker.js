addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})

const method = {
    method: "GET"
};

async function getHtml(url) {
    const response = await fetch(url, method);
    var resptxt = await response.text();
    return resptxt;
}

async function htmlhandle(request) {
    var urls = new URL(request.url);
    var base = "https://raw.githubusercontent.com/xxzl0130/GF_Tool_Server/master/HTML/";

    var json = JSON.parse(await getHtml(base + "list.json"));
    if (urls.pathname == "/") {
        var html = `<html lang="zh">
<head>
<title>轩轩的少前工具站</title>
<meta charset="utf-8">
</head>
<link rel="icon" href="https://static.xuanxuan.tech/favicon.ico" type="image/x-icon" />
<body>
<h2>轩轩的少前工具站</h2>
<p>`

        for (var item in json){
            html += '<a href=\"' + json[item].path + '\">' + json[item].name + '</a><br>\n'
        }
        html += `</p>\n</body>\n</html>`

        return html;
    }
    else{
        for (var item in json){
            if(json[item].path == urls.pathname){
                return getHtml(base + json[item].file)
            }
        }
    }
    return "404 Not Found!";
}

/**
 * Respond to the request
 * @param {Request} request
 */
async function handleRequest(request) {
    if (new URL(request.url).protocol != "https:") {
        var rhttps = new Response("Location to https", { status: 301 });
        rhttps.headers.set("Location", request.url.replace("http://", "https://"));
        return rhttps;
    }
    var resp = new Response(await htmlhandle(request), { "Content-type": "text/html;charset=UTF-8", status: 200 });
    resp.headers.set("Content-Type", "text/html");
    return resp;
}
