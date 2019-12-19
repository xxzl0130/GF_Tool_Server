addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})

var style = `<style>
	#submit {
		margin: 1rem auto;
		display: block;
		padding: .3rem 2rem;
		font-size: 1.1rem;
	}

	footer,
	h2 {
		text-align: center;
	}

	table {
		border-collapse: collapse;
	}

	#msg {
		color: #586e75;
	}

	input,
	textarea {
		width: 85%;
	}

	.responsive-label {
		align-items: center;
	}

	@media (min-width: 768px) {
		.responsive-label .col-md-3 {
			text-align: right
		}
	}

	footer {
		position: relative;
		margin-top: 3rem;
	}
</style>
`

async function htmlhandle(request){
    var urls = new URL(request.url);

    var title = "少前芯片数据查询"
    var url = "http://ar.xuanxuan.tech:8080/chip"
    var script = ""
    var body = ""
    
    if(urls.pathname == "/" || urls.pathname == "/chip"){
        title = "少前芯片数据查询"
        url = "https://ar.xuanxuan.tech:8080/chip"
        script = `button.addEventListener('click', async _ => {
			document.getElementById('msg').innerHTML = '正在查询……';

			let uid = getVal('input[name="uid"]');
			let name = getVal('input[name="name"]');
			let locked = getVal('input[name="lock-status"]:checked');
			let equipped = getVal('input[name="equipped-status"]:checked');
			let reqData = {};
			reqData.uid = uid.value;
			reqData.name = name.value;
			reqData.locked = locked.value;
			reqData.equipped = equipped.value;

			try {
				var xhr = new XMLHttpRequest();
				xhr.open("POST", url, true);
				xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
				xhr.send(serialize(reqData));

				xhr.onload = function () {
					document.getElementById('msg').innerHTML = '查询结束。';
					document.getElementById('chipResult').value = this.responseText;
				}
			} catch (err) {
                console.error('Get visit number error:' + err);
			}
		});
    };
</script>
`
        body = `<body>
        <h2>少前芯片数据查询</h2>
        <div class="container col-sm-12 col-md-8 col-lg-6">
            <form id="myForm" action="" method="POST">
                <fieldset>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="uid">UID</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="number" step=1 id="uid" name="uid" value="" required>
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="name">昵称</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="text" id="name" name="name" value="" required>
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="lock-status">输出芯片的<b>锁定</b>状态:</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="radio" name="lock-status" value=0 checked>任意
                            <input type="radio" name="lock-status" value=1>已锁定
                            <input type="radio" name="lock-status" value=2>未锁定
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="lock-status">输出芯片的<b>装备</b>状态:</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="radio" name="equipped-status" value=0>任意
                            <input type="radio" name="equipped-status" value=1>已装备
                            <input type="radio" name="equipped-status" value=2 checked>未装备
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <button type="button" class="shadowed" id="submit">查询</button>
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label>芯片代码</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <textarea autocomplete="off" id="chipResult" value=""></textarea>
                        </div>
                    </div>
                </fieldset>
            </form>
            <span id="msg"></span>
            <p><mark class="tertiary">说明</mark> 将手机连接任意WiFi，长按打开WiFi的高级设置，代理选择“自动代理”，代理地址为：</p>
            <pre>http://static.xuanxuan.tech/GF/GF.pac</pre>
            <p>然后重启游戏，登陆完成后在本页面填写UID和昵称并选择芯片配置，然后查询。
                <br>芯片代码请<code>Ctrl + A</code>全选复制。<b>数据有效期为10分钟。</b>
                <br>台服需要下载安装<a href="http://static.xuanxuan.tech/GF/ca.crt">证书</a>。
            </p>
        </div>
    
        <footer>
            <p><a href="/kalina">获取格林娜数据</a>
                | <a href="https://github.com/xxzl0130/GF_Tool_Server">GitHub</a>
                <span id="pvSpan" style="display: none;"> | 访问量: <a id="pv" href=""></a></span>
            </p>
        </footer>
        <!-- 如不需要访问统计请注释下方js -->
        <script src="https://www.fc4soda.moe/kis3/kis3.js"></script>
    </body>
`
    }
    else if(urls.pathname == "/kalina"){
        title = "格林娜好感度查询"
        url = "https://ar.xuanxuan.tech:8080/kalina"
        script = `button.addEventListener('click', async _ => {
			document.getElementById('msg').innerHTML = '正在查询……';

			let uid = getVal('input[name="uid"]');
			let name = getVal('input[name="name"]');
			let reqData = {};
			reqData.uid = uid.value;
			reqData.name = name.value;

			try {
				var xhr = new XMLHttpRequest();
				xhr.open("POST", url, true);
				xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
				xhr.send(serialize(reqData));

				xhr.onload = function () {
					document.getElementById('msg').innerHTML = '查询结束。';
					let result = this.responseText.split(';');
					document.getElementById('favorability').value = result[0];
					document.getElementById('money').value = result[1];
				}
			} catch (err) {
                console.error('Get visit number error:' + err);
			}
		});
	};
</script>
`
        body = `<body>
        <h2>格林娜好感度查询</h2>
        <div class="container col-sm-12 col-md-8 col-lg-6">
            <form id="myForm" action="" method="POST">
                <fieldset>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="uid">UID</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="number" step=1 id="uid" name="uid" value="" required>
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="name">昵称</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input type="text" id="name" name="name" value="" required>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm">
                            <button type="button" class="shadowed" id="submit">查询</button>
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="favorability">好感度</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input autocomplete="off" type="text" id="favorability" name="favorability" value="">
                        </div>
                    </div>
                    <div class="row responsive-label">
                        <div class="col-sm-12 col-md-3">
                            <label for="money">氪金量</label>
                        </div>
                        <div class="col-sm-12 col-md">
                            <input autocomplete="off" type="text" id="money" name="money" value="">
                        </div>
                    </div>
                </fieldset>
            </form>
            <span id="msg"></span>
            <p><mark class="tertiary">说明</mark> 将手机连接任意WiFi，长按打开WiFi的高级设置，代理选择“自动代理”，代理地址为：</p>
            <pre>http://static.xuanxuan.tech/GF/GF.pac</pre>
            <p>然后重启游戏，登陆完成后在本页面填写UID和昵称，然后查询。
                <br><b>数据有效期为10分钟。</b>
                <br>台服需要下载安装<a href="http://static.xuanxuan.tech/GF/ca.crt">证书</a>。
            </p>
        </div>
    
        <footer>
            <p><a href="/chip">获取芯片数据</a>
                | <a href="https://github.com/xxzl0130/GF_Tool_Server">GitHub</a>
                <span id="pvSpan" style="display: none;"> | 访问量: <a id="pv" href=""></a></span>
            </p>
        </footer>
        <!-- 如不需要访问统计请注释下方js -->
        <script src="https://www.fc4soda.moe/kis3/kis3.js"></script>
    </body>
`
    }
    else{
        return "404 Not Found!"
    }

    var header = `<html lang="zh">

<head>
	<title>${title}</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<link rel="icon" href="https://static.xuanxuan.tech/favicon.ico" type="image/x-icon" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/mini.css/3.0.1/mini-default.min.css" />
<script>
	window.onload = function () {
		var url = '${url}';
		const button = document.getElementById('submit');

		let getPv = async function () {
			const pvUrl = 'https://www.fc4soda.moe/kis3/stats?view=count&url=gf.xuanxuan.tech:8080/chip&format=json';
			const pvChartUrl = 'https://www.fc4soda.moe/kis3/stats?view=months&url=gf.xuanxuan.tech:8080/chip&format=chart';
			try {
				const response = await fetch(pvUrl, {
					method: 'get'
				}).then(response => response.json())
					.then(function (json) {
						let pv = document.getElementById('pv');
						pv.text = json[0].second;
						pv.setAttribute('href', pvChartUrl);
						document.getElementById('pvSpan').style.display = 'inline';
					});
			} catch (err) {
                console.error('Get visit number error:' + err);
			}
		}
		/** 如不需要访问统计请注释getPv调用 **/
		getPv()

		let getVal = function (selector) {
			return document.querySelector(selector);
		}

		let serialize = function (obj) {
			var str = [];
			for (var p in obj)
				if (obj.hasOwnProperty(p)) {
					str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
				}
			return str.join("&");
        }
`
    data = header + script + style + body + `</html>`
    return data
}

/**
 * Respond to the request
 * @param {Request} request
 */
async function handleRequest(request) {
    if(new URL(request.url).protocol != "https:") {
        var rhttps = new Response("Location to https", {status: 301});
		rhttps.headers.set("Location", request.url.replace("http://", "https://"));
		return rhttps;
    }
    var resp = new Response(await htmlhandle(request), {"Content-type": "text/html;charset=UTF-8", status: 200});
    resp.headers.set("Content-Type", "text/html");
    return resp;
}
