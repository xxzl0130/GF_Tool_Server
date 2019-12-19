addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})

async function htmlhandle(request) {
    var urls = new URL(request.url);
    const init = {
        method: "GET"
    };
    var base = "https://raw.githubusercontent.com/xxzl0130/GF_Tool_Server/master/HTML";
    if (urls.pathname == "/" || urls.pathname == "/chip") {
        const response = await fetch(base + "/chip.html", init);
		var resptxt = await response.text();
        return resptxt;
    }
    else if (urls.pathname == "/kalina") {
        const response = await fetch(base + "/kalina.html", init);
		var resptxt = await response.text();
        return resptxt;
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
