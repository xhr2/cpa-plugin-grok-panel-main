package main

const htmlPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Nunito:wght@500;600;700;800&family=Quicksand:wght@500;600;700&display=swap" rel="stylesheet">
<title>Grok 面板 v1.1.26</title>
<style>
:root{
--bg:#14161c;--card:#1c1f28;--card2:#252a36;--ink:#eef1f7;--muted:#9aa3b5;--line:#343b4a;--soft:#2a303c;--soft2:#353c4c;
--green:#4fd18b;--red:#ff6b7a;--yellow:#f0c14d;--blue:#6aa8ff;--violet:#a78bfa;--orange:#ff9f6b;
--green-bg:rgba(79,209,139,.14);--red-bg:rgba(255,107,122,.14);--yellow-bg:rgba(240,193,77,.14);--orange-bg:rgba(255,159,107,.14);--super-bg:rgba(167,139,250,.14);--heavy-bg:rgba(240,193,77,.14);
--accent:var(--green);--accent-bg:var(--green-bg);--focus:var(--blue);--armed-fg:#0f1410;
--radius-xs:6px;--radius-sm:8px;--radius-md:10px;--radius-lg:12px;--radius-xl:14px;--radius-pill:999px;
--font-body:"Nunito","Quicksand","Segoe UI Variable Text","Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei UI","Noto Sans SC","Source Han Sans SC",system-ui,-apple-system,BlinkMacSystemFont,sans-serif;
--font-display:"Nunito","Quicksand","Segoe UI Variable Display","Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei UI","Noto Sans SC","Source Han Sans SC",system-ui,-apple-system,BlinkMacSystemFont,sans-serif;
--font-num:ui-rounded,"SF Pro Rounded","Nunito","Segoe UI",system-ui,sans-serif;
--shadow-soft:0 .35vh .9vw rgba(0,0,0,.18);
color-scheme:dark
}
@media (prefers-color-scheme:light){
:root{
--bg:#f4f1ec;--card:#fffdfb;--card2:#f7f3ee;--ink:#2a2430;--muted:#7a7382;--line:#e5ddd4;--soft:#efe9e2;--soft2:#e7dfd6;
--green:#2f9d66;--red:#e0525e;--yellow:#c7922a;--blue:#3d7fd6;--violet:#7c5cbf;--orange:#e07a45;
--green-bg:rgba(47,157,102,.12);--red-bg:rgba(224,82,94,.11);--yellow-bg:rgba(199,146,42,.12);--orange-bg:rgba(224,122,69,.12);--super-bg:rgba(124,92,191,.11);--heavy-bg:rgba(199,146,42,.12);
--accent:var(--orange);--accent-bg:var(--orange-bg);--focus:var(--orange);--armed-fg:#fffaf6;
--shadow-soft:0 .45vh 1.1vw rgba(74,54,36,.08);
color-scheme:light
}
.brand h1{color:var(--orange)}
button:hover:not(:disabled){background:var(--orange);border-color:var(--orange);color:var(--armed-fg)}
button.armed{background:var(--orange);border-color:var(--orange);color:var(--armed-fg)}
.panel-tab.active{background:var(--orange);border-color:var(--orange);color:var(--armed-fg)}
.bar-fill,.chart-bar{background:var(--orange);border-color:var(--orange)}
.status-dot{background:var(--orange);border-color:var(--orange)}
.stat-card:hover,.panel:hover{border-color:rgba(224,122,69,.45)}
.checkline:has(input:checked){border-color:var(--orange);background:var(--orange-bg)}
input:focus,select:focus,button:focus{outline-color:var(--orange)}
input[type=checkbox]:checked{background:var(--orange);border-color:var(--orange)}
}
*{box-sizing:border-box}
html,body{min-height:100vh}
body{
margin:0;background:
radial-gradient(circle at 12% -8%,rgba(106,168,255,.12),transparent 34vw),
radial-gradient(circle at 92% 0%,rgba(167,139,250,.1),transparent 30vw),
var(--bg);
color:var(--ink);font-family:var(--font-body);
/* Dense panel scale: clamp prevents giant UI on wide monitors / management iframe */
font-size:clamp(12px,.72vw,14px);line-height:1.45;
font-weight:500;letter-spacing:.01em;padding:clamp(10px,.9vw,16px);
-webkit-font-smoothing:antialiased;text-rendering:optimizeLegibility
}
button,input,select{font-family:inherit;color:var(--ink)}
button{
background:var(--card2);border:.1vw solid var(--line);padding:6px 12px;min-height:32px;cursor:pointer;
font-size:clamp(12px,.72vw,13.5px);font-weight:700;text-transform:none;letter-spacing:.01vw;border-radius:var(--radius-pill);
box-shadow:var(--shadow-soft)
}
button:hover:not(:disabled){background:var(--soft2);border-color:var(--muted);color:var(--ink)}
button.armed{background:var(--red);border-color:var(--red);color:#fff}
button:disabled{cursor:not-allowed;opacity:.42;background:var(--soft);box-shadow:none}
input,select{
background:var(--card2);border:.1vw solid var(--line);padding:6px 10px;min-height:32px;font-size:clamp(12px,.72vw,13.5px);
border-radius:var(--radius-md);color:var(--ink);font-weight:600
}
input:focus,select:focus,button:focus{outline:.16vw solid color-mix(in srgb,var(--focus) 55%,transparent);outline-offset:.12vw}
input[type=checkbox]{
appearance:none;width:1.15vw;height:1.15vw;min-height:1.15vw;padding:0;margin:0;background:var(--card2);
border:.1vw solid var(--muted);vertical-align:middle;cursor:pointer;border-radius:.4vw;position:relative
}
input[type=checkbox]:checked{background:var(--accent);border-color:var(--accent)}
input[type=checkbox]:checked::after{
content:"";position:absolute;left:50%;top:45%;width:.28vw;height:.52vw;border:solid #fff;border-width:0 .14vw .14vw 0;
transform:translate(-50%,-55%) rotate(45deg)
}
.shell{width:min(96vw,1280px);margin:0 auto}
.topline{display:flex;justify-content:space-between;align-items:center;gap:10px;margin-bottom:10px}
.brand{display:flex;align-items:center;gap:.75vw;min-height:34px}
h1{
font-family:var(--font-display);font-size:clamp(18px,1.2vw,22px);line-height:1.05;margin:0;font-weight:800;
letter-spacing:-.02vw;text-transform:none
}
h2{
font-family:var(--font-display);font-size:clamp(14px,.9vw,16px);margin:0 0 1vh 0;font-weight:800;letter-spacing:.01vw;
border-bottom:.1vw solid var(--line);padding-bottom:.75vh;text-transform:none
}
.status-dot{width:.95vw;height:.95vw;background:var(--accent);display:inline-block;border:.08vw solid color-mix(in srgb,var(--accent) 70%,#fff);border-radius:50%;box-shadow:0 0 .45vw color-mix(in srgb,var(--accent) 45%,transparent)}
.status-dot.err{background:var(--red);border-color:var(--red);box-shadow:0 0 .45vw rgba(255,107,122,.4)}
.muted{color:var(--muted);font-size:clamp(11.5px,.68vw,13px);font-weight:600}
.top-actions{display:flex;gap:.65vw;align-items:center;flex-wrap:wrap;justify-content:flex-end}
.feedback{
border:.1vw solid var(--line);background:var(--card);padding:8px 12px;min-height:34px;margin-bottom:1.1vh;
color:var(--muted);font-size:clamp(12px,.72vw,13.5px);border-radius:var(--radius-lg);font-weight:600
}
.feedback.ok{border-color:color-mix(in srgb,var(--green) 55%,var(--line));color:var(--green);background:var(--green-bg)}
.feedback.warn{border-color:color-mix(in srgb,var(--yellow) 55%,var(--line));color:var(--yellow);background:var(--yellow-bg)}
.feedback.error{border-color:color-mix(in srgb,var(--red) 55%,var(--line));color:var(--red);background:var(--red-bg)}
.stats-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:10px;margin-bottom:12px}
.stat-card,.panel{
background:var(--card);border:.1vw solid var(--line);padding:10px 12px;border-radius:var(--radius-xl);
box-shadow:var(--shadow-soft)
}
.stat-card{min-height:96px;display:flex;flex-direction:column;justify-content:space-between}
.stat-card:hover,.panel:hover{border-color:color-mix(in srgb,var(--accent) 40%,var(--line))}
.stat-label{font-size:clamp(11px,.62vw,12px);color:var(--muted);text-transform:none;letter-spacing:.02vw;font-weight:700}
.stat-value{font-family:var(--font-num);font-size:clamp(18px,1.2vw,24px);line-height:1.08;font-weight:800;margin-top:.65vh;word-break:break-word;letter-spacing:-.02vw}
.stat-sub{font-size:clamp(11px,.62vw,12px);color:var(--muted);margin-top:.55vh;font-weight:600}
.bar-container{width:100%;max-width:18vw;height:1.15vh;background:var(--soft2);margin-top:.75vh;border:.06vw solid var(--line);overflow:hidden;border-radius:var(--radius-pill)}
.bar-container.small{max-width:11vw;height:1vh;margin-top:.4vh}
.bar-fill{height:100%;background:var(--accent);border-radius:var(--radius-pill)}
.bar-fill.warn{background:var(--yellow)}
.bar-fill.danger{background:var(--red)}
.panel{margin-bottom:12px}
.form-grid{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:.85vw;align-items:end}
.field{display:flex;flex-direction:column;gap:.48vh}
.field label,.checkline{font-size:clamp(11.5px,.66vw,12.5px);color:var(--muted);font-weight:650}
.number-input{width:100%}
.checkline{
display:flex;align-items:center;gap:.5vw;min-height:32px;border:.1vw solid var(--soft2);padding:.65vh .75vw;
background:var(--bg);border-radius:var(--radius-md)
}
.checkline:has(input:checked){border-color:var(--accent);background:var(--accent-bg)}
.help-text{font-size:clamp(11px,.62vw,12px);color:var(--muted);margin-top:.9vh;font-weight:600}
.field-wide{grid-column:1 / -1}
.auth-grid{margin-top:1vh}
.auth-badge{margin-top:1vh;padding:.85vh .9vw;border:.1vw solid var(--soft2);font-size:clamp(11.5px,.66vw,12.5px);background:var(--bg);border-radius:var(--radius-md);font-weight:650}
.auth-badge.ok{border-color:var(--green);color:var(--green)}
.auth-badge.warn{border-color:var(--orange);color:var(--orange)}
.chart-row{
display:flex;align-items:flex-end;gap:.28vw;height:16vh;padding:1vh .7vw;overflow-x:auto;overflow-y:hidden;
background:var(--bg);border:.1vw solid var(--line);border-radius:var(--radius-lg)
}
.chart-bar{flex:0 0 1.25vw;min-height:.35vh;background:var(--accent);opacity:.78;position:relative;border:0;border-radius:.45vw .45vw .2vw .2vw}
.chart-bar.fail{background:var(--red)}
.chart-bar:hover{opacity:1}
.chart-bar:hover::after{
content:attr(data-tip);position:absolute;bottom:14.5vh;left:0;background:var(--card2);color:var(--ink);
padding:.5vh .6vw;font-size:clamp(10.5px,.58vw,11.5px);white-space:normal;z-index:8;width:14vw;border:.08vw solid var(--line);
border-radius:var(--radius-sm);box-shadow:var(--shadow-soft);font-weight:650
}
.chart-empty{color:var(--muted);padding:5vh 1vw;font-size:.82vw;font-weight:650}
.filter-grid{display:grid;grid-template-columns:minmax(0,1.4fr) repeat(4,minmax(0,.75fr)) minmax(0,1fr) minmax(0,.7fr);gap:.75vw;margin-bottom:1vh;align-items:end}
.search-box,.select-filter,.sort-filter{width:100%}
.batchbar,.pagebar{
display:flex;align-items:center;gap:.7vw;flex-wrap:wrap;margin-bottom:.95vh;background:var(--card2);
border:.1vw solid var(--line);padding:.85vh .9vw;border-radius:var(--radius-lg)
}
.batchbar .spacer{flex:1}
.pagebar{margin:.95vh 0}
.pagebar .page-info{font-size:clamp(11.5px,.66vw,12.5px);color:var(--muted);font-weight:650}
.pagebar .page-info b{color:var(--ink)}
.pagebar .field{margin:0}
.pagebar select{min-width:8vw}
.table-wrap{max-height:min(62vh,640px);overflow:auto;border:.1vw solid var(--line);background:var(--card);border-radius:var(--radius-xl);width:100%}
table{border-collapse:separate;border-spacing:0;width:100%;min-width:0;table-layout:fixed;background:transparent}
th,td{text-align:left;padding:.7vh .55vw;border-bottom:.08vw solid var(--line);font-size:clamp(11.5px,.66vw,12.5px);vertical-align:middle}
th{
position:sticky;top:0;z-index:4;background:var(--card2);color:var(--ink);font-family:var(--font-display);
font-weight:800;text-transform:none;font-size:clamp(10.5px,.6vw,12px);letter-spacing:.01vw;border-bottom:.12vw solid var(--line);white-space:nowrap
}
tr:hover{background:color-mix(in srgb,var(--card2) 78%,transparent)}
.col-check{width:3.2vw;text-align:center}
.col-idx{width:3.2vw;color:var(--muted);font-variant-numeric:tabular-nums}
.col-type{width:7vw}
.col-health{width:9vw}
.col-usage{width:16vw}
.col-actions{width:14vw}
.email-cell{width:auto;min-width:0;word-break:break-word;overflow-wrap:anywhere}
.email-main{font-weight:700;line-height:1.25}
.email-meta{display:flex;flex-wrap:wrap;gap:.35vw;margin-top:.35vh;align-items:center}
.num-cell{font-variant-numeric:tabular-nums;white-space:nowrap}
.usage-cell .metric-compact{display:flex;justify-content:space-between;gap:.4vw;font-size:clamp(10.5px,.58vw,11.5px);color:var(--muted);font-weight:650;margin-bottom:.3vh}
.usage-cell .metric-compact b{color:var(--ink);font-family:var(--font-num)}
.usage-cell .bar-container.small{width:100%;max-width:100%;height:.85vh;margin-top:.2vh}
.actions-cell{display:flex;gap:.35vw;flex-wrap:wrap;align-items:center}
.actions-cell button{padding:.45vh .55vw;min-height:3.1vh;font-size:clamp(10.5px,.6vw,12px);border-radius:var(--radius-pill)}
.tag{
display:inline-block;padding:.35vh .65vw;border:.08vw solid;font-size:clamp(10.5px,.58vw,11.5px);font-weight:750;letter-spacing:.01vw;
text-transform:none;min-width:4.6vw;text-align:center;border-radius:var(--radius-pill)
}
.tag.active,.tag.healthy,.tag.free{background:var(--green-bg);border-color:color-mix(in srgb,var(--green) 55%,var(--line));color:var(--green)}
.tag.disabled,.tag.invalid{background:var(--red-bg);border-color:color-mix(in srgb,var(--red) 55%,var(--line));color:var(--red)}
.tag.warn{background:var(--yellow-bg);border-color:color-mix(in srgb,var(--yellow) 55%,var(--line));color:var(--yellow)}
.tag.unknown{background:var(--soft);border-color:var(--muted);color:var(--muted)}
.tag.super{background:var(--super-bg);border-color:color-mix(in srgb,var(--violet) 55%,var(--line));color:var(--violet)}
.tag.heavy{background:var(--heavy-bg);border-color:color-mix(in srgb,var(--yellow) 55%,var(--line));color:var(--yellow)}
.tag.other{background:var(--soft);border-color:var(--muted);color:var(--muted)}
.tag.error,.tag.unavailable{background:var(--red-bg);border-color:color-mix(in srgb,var(--red) 55%,var(--line));color:var(--red)}
.health-dot.unavailable{background:var(--red)}
.cell-sub{font-size:clamp(10px,.56vw,11px);color:var(--muted);margin-top:.35vh;word-break:break-word;font-weight:600}
.red-text{color:var(--red);font-weight:750}
.summary-row{display:flex;gap:1.2vw;flex-wrap:wrap;font-size:clamp(11.5px,.66vw,12.5px);color:var(--muted);margin-top:.9vh;font-weight:650}
.summary-row b{color:var(--ink)}
.row-invalid{background:var(--red-bg)!important}
.row-disabled{background:var(--soft)!important;opacity:.62}
.row-warn{background:var(--yellow-bg)!important}
.health-indicator{display:inline-flex;align-items:center;gap:.4vw}
.health-dot{width:.72vw;height:.72vw;border-radius:50%;display:inline-block;flex-shrink:0}
.health-dot.healthy{background:var(--green);box-shadow:0 0 .4vw color-mix(in srgb,var(--green) 55%,transparent)}
.health-dot.invalid{background:var(--red);box-shadow:0 0 .4vw color-mix(in srgb,var(--red) 55%,transparent)}
.health-dot.warn{background:var(--yellow);box-shadow:0 0 .4vw color-mix(in srgb,var(--yellow) 55%,transparent)}
.health-dot.disabled{background:var(--muted)}
.health-dot.unknown{background:var(--soft2);border:.08vw solid var(--muted)}
.panel-tabs{
display:flex;gap:.55vw;margin-bottom:1.4vh;border:.1vw solid var(--line);background:var(--card);
padding:.55vh .55vw;position:sticky;top:.4vh;z-index:20;border-radius:var(--radius-xl);box-shadow:var(--shadow-soft)
}
.panel-tab{flex:1;background:transparent;color:var(--muted);font-weight:750;border-color:transparent;box-shadow:none}
.panel-tab:hover:not(:disabled){background:var(--soft);color:var(--ink)}
.panel-tab.active{background:var(--accent);border-color:var(--accent);color:var(--armed-fg);box-shadow:none}
.panel-view{display:none}
.panel-view.active{display:block}
.panel-view>.panel:last-child{margin-bottom:0}
.metric-line{display:flex;justify-content:space-between;gap:.5vw;align-items:baseline}
.metric-line span{color:var(--muted);font-size:clamp(10.5px,.6vw,12px);font-weight:650}
.metric-line b{color:var(--ink);font-size:clamp(12px,.72vw,13.5px);font-family:var(--font-num)}
.meter-head{display:flex;justify-content:space-between;gap:.5vw;font-size:clamp(10.5px,.58vw,11.5px);color:var(--muted);font-weight:650}
@media (prefers-color-scheme:dark){
body{background:
radial-gradient(circle at 10% -10%,rgba(106,168,255,.16),transparent 32vw),
radial-gradient(circle at 95% 4%,rgba(167,139,250,.12),transparent 28vw),
var(--bg)}
}
@media (orientation:portrait),(hover:none) and (pointer:coarse){
body{font-size:3.2vw;padding:3vw}
.shell{width:94vw}
.topline{flex-direction:column;align-items:stretch;gap:1vh}
.brand{align-items:flex-start;flex-direction:column;gap:.8vh}
.top-actions{display:grid;grid-template-columns:1fr}
.top-actions button,.batchbar button{width:100%}
h1,.brand h1{font-size:6.1vw}
h2{font-size:3.7vw}
.muted{font-size:2.8vw}
.status-dot{width:2.8vw;height:2.8vw}
.feedback{font-size:2.8vw;border-radius:3.2vw;padding:1vh 3vw}
.stats-grid,.form-grid,.filter-grid{grid-template-columns:1fr;width:100%;gap:1vh}
.stat-card,.panel{border-radius:4vw;padding:1.35vh 3vw}
.stat-card{min-height:12vh}
.stat-label,.stat-sub,.field label,.checkline,.help-text,.summary-row{font-size:2.6vw}
.stat-value{font-size:5.4vw}
.bar-container{max-width:100%;height:1.2vh}
.bar-container.small{max-width:100%;height:1vh}
.number-input,.search-box,.select-filter,.sort-filter{width:100%}
button,input,select{font-size:2.8vw;min-height:5.2vh;border-radius:2.8vw;padding:.9vh 2.6vw}
input[type=checkbox]{width:4vw;height:4vw;min-height:4vw;border-radius:1.1vw}
input[type=checkbox]:checked::after{width:.9vw;height:1.7vw;border-width:0 .45vw .45vw 0}
.checkline{gap:2vw;padding:1vh 2.4vw;border-radius:2.4vw}
.batchbar,.pagebar{border-radius:3.2vw;gap:1vh}
.pagebar{align-items:stretch}
.pagebar .page-info{font-size:2.6vw}
.chart-row{height:18vh;border-radius:3.2vw}
.chart-bar{flex-basis:3vw;border-radius:1vw 1vw .4vw .4vw}
.chart-bar:hover::after{font-size:2.4vw;width:44vw;bottom:16vh}
.table-wrap{max-height:55vh;border-radius:3.2vw}
table{width:100%;min-width:0}
th,td{font-size:2.6vw;padding:.95vh 1.4vw}
th{font-size:2.25vw}
.tag{font-size:2.25vw;min-width:12vw;padding:.4vh 1.2vw}
.cell-sub,.metric-line span,.meter-head{font-size:2.2vw}
.metric-line b{font-size:3vw}
.actions-cell button{font-size:2.3vw;min-height:4.8vh;padding:.65vh 1.5vw}.col-check{width:8vw}.col-idx{width:8vw}.col-type{width:16vw}.col-health{width:18vw}.col-usage{width:24vw}.col-actions{width:28vw}.actions-cell button{font-size:2.2vw;min-height:4.4vh;padding:.55vh 1.1vw}.email-main{font-size:2.7vw}.usage-cell .metric-compact{font-size:2.1vw}
.panel-tabs{top:.3vh;gap:1vw;padding:.7vh 1vw;border-radius:3.4vw}
.panel-tab{font-size:2.7vw;min-height:5.2vh;padding:.8vh 1vw}
}
@supports not (color:color-mix(in srgb,red,blue)){
input:focus,select:focus,button:focus{outline-color:var(--focus)}
.feedback.ok,.tag.active,.tag.healthy,.tag.free{border-color:var(--green)}
.feedback.warn,.tag.warn,.tag.heavy{border-color:var(--yellow)}
.feedback.error,.tag.disabled,.tag.invalid,.tag.error,.tag.unavailable{border-color:var(--red)}
.tag.super{border-color:var(--violet)}
.stat-card:hover,.panel:hover{border-color:var(--muted)}
.health-dot.healthy{box-shadow:0 0 .4vw var(--green)}
.health-dot.invalid{box-shadow:0 0 .4vw var(--red)}
.health-dot.warn{box-shadow:0 0 .4vw var(--yellow)}
}
</style>
</head>
<body>
<div class="shell">
<header class="topline">
<div class="brand"><h1>Grok Panel</h1><span class="status-dot" id="statusDot"></span><span class="muted" id="lastUpdate">等待数据</span></div>
<div class="top-actions"><button id="refreshBtn">刷新数据</button><button id="checkVisibleTopBtn">手动检查本页</button><button id="verifyVisibleTopBtn">核实本页套餐</button></div>
</header>
<div id="feedback" class="feedback">就绪：所有操作使用同源插件端点，不包含管理密钥。</div>
<nav class="panel-tabs" aria-label="面板导航">
<button type="button" class="panel-tab active" data-panel="stats">数据统计</button>
<button type="button" class="panel-tab" data-panel="accounts">账号管理</button>
<button type="button" class="panel-tab" data-panel="settings">设置</button>
</nav>
<div class="panel-view active" data-view="stats">
<section class="stats-grid">
<div class="stat-card"><div class="stat-label">Grok 文件总数</div><div class="stat-value" id="statTotal">--</div><div class="stat-sub" id="statTotalSub"></div></div>
<div class="stat-card"><div class="stat-label">活跃率</div><div class="stat-value" id="statActive">--</div><div class="stat-sub" id="statActiveSub"></div><div class="bar-container"><div class="bar-fill" id="activeBar"></div></div></div>
<div class="stat-card"><div class="stat-label">请求成功率</div><div class="stat-value" id="statRequests">--</div><div class="stat-sub" id="statRequestsSub"></div><div class="bar-container"><div class="bar-fill" id="requestBar"></div></div></div>
<div class="stat-card"><div class="stat-label">Token 使用率</div><div class="stat-value" id="statUsage">--</div><div class="stat-sub" id="statUsageSub"></div><div class="bar-container"><div class="bar-fill" id="usageBar"></div></div></div>
<div class="stat-card"><div class="stat-label">估算 Token</div><div class="stat-value" id="statTokens">--</div><div class="stat-sub" id="statTokensSub"></div></div>
<div class="stat-card"><div class="stat-label">总容量</div><div class="stat-value" id="statCapacity">--</div><div class="stat-sub" id="statCapacitySub"></div></div>
<div class="stat-card"><div class="stat-label">账号类型</div><div class="stat-value" id="statTypes">--</div><div class="stat-sub" id="statTypesSub"></div></div>
<div class="stat-card"><div class="stat-label">健康率</div><div class="stat-value" id="statHealth">--</div><div class="stat-sub" id="statHealthSub"></div><div class="bar-container"><div class="bar-fill" id="healthBar"></div></div></div>
<div class="stat-card"><div class="stat-label">异常 / 可安全清理</div><div class="stat-value" id="statInvalid">--</div><div class="stat-sub" id="statInvalidSub"></div></div>
<div class="stat-card"><div class="stat-label">已选择</div><div class="stat-value" id="statSelected">0</div><div class="stat-sub" id="statSelectedSub"></div></div>
</section>
<section class="panel"><h2>请求趋势</h2><div class="chart-row" id="chartRow"></div></section>
</div>
<div class="panel-view" data-view="settings">
<section class="panel">
<h2>设置</h2>
<div class="form-grid">
<div class="field"><label for="tokenLimit">每账号上限 Token</label><input class="number-input" type="number" id="tokenLimit" value="2000000" min="1"></div>
<div class="field"><label for="tokensPerReq">每请求估算 Token</label><input class="number-input" type="number" id="tokensPerReq" value="5000" min="1"></div>
<div class="field"><label for="failThreshold">失败阈值</label><input class="number-input" type="number" id="failThreshold" value="3" min="1"></div>
<label class="checkline"><input type="checkbox" id="autoCheck"> 自动检查 默认关</label>
<label class="checkline"><input type="checkbox" id="autoDelete"> 自动删除无效 默认关</label>
<label class="checkline"><input type="checkbox" id="protectSuper" checked> 保护 super 默认开</label>
<label class="checkline"><input type="checkbox" id="protectHeavy" checked> 保护 heavy 默认开</label>
<label class="checkline"><input type="checkbox" id="protectUnknown"> 保护未知 默认关（未知按 Free）</label>
</div>
<div class="help-text">删除和清理需要再次点击同一按钮确认。插件在管理中心 iframe 中通常读不到“记住密码”，请在下方填写管理密钥后保存。</div>
<div class="form-grid auth-grid">
<div class="field field-wide"><label for="mgmtKey">CPA 管理密钥（仅本机浏览器）</label><input type="password" id="mgmtKey" placeholder="粘贴 management key" autocomplete="off"></div>
<label class="checkline"><input type="checkbox" id="rememberMgmtKey" checked> 记住管理密钥 默认开</label>
<button type="button" id="saveMgmtKeyBtn">保存密钥</button>
<button type="button" id="clearMgmtKeyBtn">清除密钥</button>
</div>
<div class="auth-badge warn" id="authBadge">管理授权：检测中...</div>
</section>
</div>
<div class="panel-view" data-view="accounts">
<section class="panel"><h2>账号明细</h2>
<div class="filter-grid">
<div class="field"><label for="searchBox">搜索</label><input type="text" class="search-box" id="searchBox" placeholder="邮箱、状态、类型"></div>
<div class="field"><label for="statusFilter">状态</label><select class="select-filter" id="statusFilter"><option value="all">全部</option><option value="active">活跃</option><option value="disabled">禁用</option><option value="error">错误</option><option value="unavailable">异常</option><option value="other">其他</option><option value="unknown">未知</option></select></div>
<div class="field"><label for="typeFilter">类型</label><select class="select-filter" id="typeFilter"><option value="all">全部</option><option value="free">Free</option><option value="super">Super</option><option value="heavy">Heavy</option></select></div>
<div class="field"><label for="healthFilter">健康</label><select class="select-filter" id="healthFilter"><option value="all">全部</option><option value="healthy">健康</option><option value="warn">警告</option><option value="unavailable">异常</option><option value="invalid">无效</option><option value="disabled">禁用</option><option value="unknown">未知</option></select></div>
<div class="field"><label for="usageFilter">用量</label><select class="select-filter" id="usageFilter"><option value="all">全部</option><option value="unused">未使用</option><option value="low">低于一半</option><option value="warn">一半以上</option><option value="high">高于八成</option></select></div>
<div class="field"><label for="sortFilter">排序</label><select class="sort-filter" id="sortFilter"><option value="success_desc">成功请求降序</option><option value="failed_desc">失败降序</option><option value="usage_desc">用量降序</option><option value="health_asc">健康优先</option><option value="type_asc">类型</option><option value="email_asc">邮箱</option></select></div>
<div class="field"><label for="pageSize">每页</label><select class="select-filter" id="pageSize"><option value="20">20</option><option value="50" selected>50</option><option value="100">100</option><option value="200">200</option></select></div>
</div>
<div class="batchbar">
<label class="checkline"><input type="checkbox" id="selectVisible"> 选择本页</label>
<button id="checkSelectedBtn">检查选中</button>
<button id="verifySelectedBtn">核实选中套餐</button>
<button id="batchDeleteBtn">删除选中</button>
<button id="cleanupInvalidBtn">安全清理 0</button>
<span class="muted" id="selectionInfo">未选择</span><span class="spacer"></span><span class="muted" id="filterInfo">无过滤</span>
</div>
<div class="pagebar">
<button type="button" id="pageFirstBtn">首页</button>
<button type="button" id="pagePrevBtn">上一页</button>
<span class="page-info" id="pageInfo">第 <b>1</b> / 1 页</span>
<button type="button" id="pageNextBtn">下一页</button>
<button type="button" id="pageLastBtn">末页</button>
<span class="spacer"></span>
<span class="page-info" id="pageRange">显示 0-0 / 0</span>
</div>
<div class="table-wrap"><table><thead><tr><th class="col-check">选</th><th class="col-idx">#</th><th>账号</th><th class="col-type">类型</th><th class="col-health">健康</th><th class="col-usage">用量</th><th class="col-actions">操作</th></tr></thead><tbody id="tableBody"></tbody></table></div>
<div class="pagebar">
<button type="button" id="pagePrevBtn2">上一页</button>
<span class="page-info" id="pageInfo2">第 <b>1</b> / 1 页</span>
<button type="button" id="pageNextBtn2">下一页</button>
</div>
<div class="summary-row"><span>可见 <b id="sumCount">0</b></span><span>已使用 <b id="sumUsed">0</b></span><span>剩余 <b id="sumAvail">0</b></span><span>成功率 <b id="sumRate">0%</b></span><span>无效 <b id="sumInvalid">0</b></span><span>可清理 <b id="sumCleanup">0</b></span></div>
</section>
</div>
</div>
<script>
/*
Frontend v1.1.26 same-origin endpoint contract for a matching backend.
Delete/check reuse CPA management auth. Key resolution order:
1) panel-local saved management key
2) parent/local cli-proxy-auth (remember password)
3) same-origin fallback with empty key is rejected
GET  ./data                         -> stats + files
POST /v0/management/plugins/grok-panel/checks
DELETE /v0/management/auth-files
*/
var settingsKey='grok-panel-v1.1.20-settings';
var mgmtKeyStore='grok-panel-mgmt-key';
var allData=[];
var lastData=null;
var selected={};
var confirmUntil={};
var healthOverrides={};
var typeOverrides={};
var busy=false;
var autoCheckBusy=false;
var lastAutoCheckAt=0;
var refreshHandle=null;
var currentPage=1;
var settings=loadSettings();
function byId(id){return document.getElementById(id)}
function defaults(){return{tokenLimit:2000000,tokensPerReq:5000,threshold:3,autoCheck:false,autoDelete:false,protectSuper:true,protectHeavy:true,protectUnknown:false,rememberMgmtKey:true,pageSize:50}}
function loadSettings(){var base=defaults();try{var raw=localStorage.getItem(settingsKey);if(raw){var parsed=JSON.parse(raw);Object.keys(base).forEach(function(k){if(parsed[k]!==undefined)base[k]=parsed[k]})}}catch(e){}return base}
function saveSettings(){try{localStorage.setItem(settingsKey,JSON.stringify(settings))}catch(e){}}
function initSettings(){byId('tokenLimit').value=settings.tokenLimit;byId('tokensPerReq').value=settings.tokensPerReq;byId('failThreshold').value=settings.threshold;byId('autoCheck').checked=!!settings.autoCheck;byId('autoDelete').checked=!!settings.autoDelete;byId('protectSuper').checked=!!settings.protectSuper;byId('protectHeavy').checked=!!settings.protectHeavy;byId('protectUnknown').checked=!!settings.protectUnknown;byId('rememberMgmtKey').checked=!!settings.rememberMgmtKey;if(byId('pageSize'))byId('pageSize').value=String([20,50,100,200].indexOf(Number(settings.pageSize))>=0?Number(settings.pageSize):50);var saved=loadManualManagementKey();if(saved)byId('mgmtKey').value=saved;updateAuthBadge()}
function readSettings(evt){settings.tokenLimit=Math.max(1,parseInt(byId('tokenLimit').value,10)||2000000);settings.tokensPerReq=Math.max(1,parseInt(byId('tokensPerReq').value,10)||5000);settings.threshold=Math.max(1,parseInt(byId('failThreshold').value,10)||3);settings.autoCheck=!!byId('autoCheck').checked;settings.autoDelete=!!byId('autoDelete').checked;settings.protectSuper=!!byId('protectSuper').checked;settings.protectHeavy=!!byId('protectHeavy').checked;settings.protectUnknown=!!byId('protectUnknown').checked;settings.rememberMgmtKey=!!byId('rememberMgmtKey').checked;settings.pageSize=Math.max(1,parseInt(byId('pageSize')&&byId('pageSize').value,10)||50);saveSettings();if(evt&&evt.target&&evt.target.id==='autoDelete'&&settings.autoDelete)setFeedback('自动删除已开启：只处理未受保护且判定无效的账号。','warn');if(evt&&evt.target&&evt.target.id==='autoCheck'&&settings.autoCheck)setFeedback('自动检查已开启：刷新后会调用同源检查端点。','warn');renderAll();updateAuthBadge()}
function switchPanel(name){name=['stats','accounts','settings'].indexOf(name)>=0?name:'stats';document.querySelectorAll('.panel-tab').forEach(function(btn){var active=btn.getAttribute('data-panel')===name;btn.classList.toggle('active',active);btn.setAttribute('aria-selected',active?'true':'false')});document.querySelectorAll('.panel-view').forEach(function(view){view.classList.toggle('active',view.getAttribute('data-view')===name)});try{sessionStorage.setItem('grok-panel-active-view',name)}catch(e){}if(name==='accounts')window.setTimeout(function(){renderTable()},0)}
function bindEvents(){document.querySelectorAll('.panel-tab').forEach(function(btn){btn.addEventListener('click',function(){switchPanel(btn.getAttribute('data-panel'))})});byId('refreshBtn').addEventListener('click',function(){fetchData(true,false)});byId('checkVisibleTopBtn').addEventListener('click',manualCheckVisible);byId('verifyVisibleTopBtn').addEventListener('click',verifyVisibleTiers);byId('checkSelectedBtn').addEventListener('click',manualCheckSelected);byId('verifySelectedBtn').addEventListener('click',verifySelectedTiers);byId('batchDeleteBtn').addEventListener('click',requestBatchDelete);byId('cleanupInvalidBtn').addEventListener('click',requestCleanupInvalid);byId('selectVisible').addEventListener('change',toggleVisibleSelection);byId('saveMgmtKeyBtn').addEventListener('click',saveManualManagementKey);byId('clearMgmtKeyBtn').addEventListener('click',clearManualManagementKey);['searchBox','statusFilter','typeFilter','healthFilter','usageFilter','sortFilter'].forEach(function(id){byId(id).addEventListener('input',function(){currentPage=1;renderAll()});byId(id).addEventListener('change',function(){currentPage=1;renderAll()})});if(byId('pageSize'))byId('pageSize').addEventListener('change',function(){settings.pageSize=Math.max(1,parseInt(byId('pageSize').value,10)||50);saveSettings();currentPage=1;renderAll()});function goPage(delta,abs){var meta=getPageMeta(getFilteredData());if(abs!==undefined)currentPage=abs;else currentPage+=delta;currentPage=Math.max(1,Math.min(meta.totalPages,currentPage));renderAll()}['pageFirstBtn','pagePrevBtn','pageNextBtn','pageLastBtn','pagePrevBtn2','pageNextBtn2'].forEach(function(id){var el=byId(id);if(!el)return;el.addEventListener('click',function(){if(id==='pageFirstBtn')goPage(0,1);else if(id==='pagePrevBtn'||id==='pagePrevBtn2')goPage(-1);else if(id==='pageNextBtn'||id==='pageNextBtn2')goPage(1);else if(id==='pageLastBtn'){var meta=getPageMeta(getFilteredData());goPage(0,meta.totalPages)}})});['tokenLimit','tokensPerReq','failThreshold','autoCheck','autoDelete','protectSuper','protectHeavy','protectUnknown','rememberMgmtKey'].forEach(function(id){byId(id).addEventListener('input',readSettings);byId(id).addEventListener('change',readSettings)});byId('tableBody').addEventListener('click',handleTableClick);byId('tableBody').addEventListener('change',handleTableChange);window.addEventListener('resize',function(){renderAll()})}
function fmt(n){n=Number(n)||0;if(n>=1000000000)return(n/1000000000).toFixed(2)+'B';if(n>=1000000)return(n/1000000).toFixed(2)+'M';if(n>=1000)return(n/1000).toFixed(1)+'K';return String(n)}
function fmtTime(){var d=new Date();return d.toLocaleTimeString('zh-CN',{hour12:false})+' '+d.toLocaleDateString('zh-CN')}
function esc(v){return String(v===undefined||v===null?'':v).replace(/[&<>"']/g,function(c){return{'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]})}
function apiBase(){return String(window.location.pathname||'').replace(/\/+$/,'')}
function fixedApiUrl(endpoint){return apiBase()+'/'+String(endpoint).replace(/^\/+/, '')}
function xorDecode(bytes,key){var out=new Uint8Array(bytes.length);for(var i=0;i<bytes.length;i++)out[i]=bytes[i]^key[i%key.length];return new TextDecoder().decode(out)}
function storageGet(key){var stores=[];try{stores.push(localStorage)}catch(e){}try{if(window.parent&&window.parent!==window)stores.push(window.parent.localStorage)}catch(e){}try{stores.push(sessionStorage)}catch(e){}for(var i=0;i<stores.length;i++){try{var v=stores[i].getItem(key);if(v)return v}catch(e){}}return null}
function storageSet(key,value,useLocal){try{if(useLocal)localStorage.setItem(key,value);else sessionStorage.setItem(key,value);return true}catch(e){return false}}
function storageRemove(key){try{localStorage.removeItem(key)}catch(e){}try{sessionStorage.removeItem(key)}catch(e){}}
function decodeCliProxyAuth(raw){if(!raw)return null;var text=String(raw);if(text.indexOf('enc::v1::')===0){var bin=atob(text.slice(9));var bytes=new Uint8Array(bin.length);for(var i=0;i<bin.length;i++)bytes[i]=bin.charCodeAt(i);var hosts=[window.location.host];try{if(window.parent&&window.parent!==window&&window.parent.location)hosts.push(window.parent.location.host)}catch(e){}var uas=[navigator.userAgent];try{if(window.parent&&window.parent.navigator)uas.push(window.parent.navigator.userAgent)}catch(e){}var decoded=null;hosts.forEach(function(host){uas.forEach(function(ua){if(decoded)return;try{decoded=xorDecode(bytes,new TextEncoder().encode('cli-proxy-api-webui::secure-storage|'+host+'|'+ua))}catch(e){}})});if(!decoded)decoded=xorDecode(bytes,new TextEncoder().encode('cli-proxy-api-webui::secure-storage'));text=decoded}try{var parsed=JSON.parse(text);if(typeof parsed==='string'){try{parsed=JSON.parse(parsed)}catch(e){}}var state=parsed&&parsed.state?parsed.state:parsed;if(state&&typeof state==='object')return state;return null}catch(e){return null}}
function loadManualManagementKey(){return String(storageGet(mgmtKeyStore)||'').trim()}
function saveManualManagementKey(){var key=String(byId('mgmtKey').value||'').trim();if(!key){setFeedback('请先输入管理密钥。','warn');return}storageSet(mgmtKeyStore,key,!!settings.rememberMgmtKey);if(settings.rememberMgmtKey){try{sessionStorage.removeItem(mgmtKeyStore)}catch(e){}}else{try{localStorage.removeItem(mgmtKeyStore)}catch(e){}}setFeedback('管理密钥已保存到本面板（仅当前浏览器）。删除/检查已可用。','ok');updateAuthBadge()}
function clearManualManagementKey(){storageRemove(mgmtKeyStore);byId('mgmtKey').value='';setFeedback('已清除本面板保存的管理密钥。','warn');updateAuthBadge()}
function connectionFromState(state,source){if(!state)return null;var key=String(state.managementKey||state.management_key||state.key||'').trim();if(!key)return null;return{apiBase:String(state.apiBase||state.api_base||state.apiUrl||'').replace(/\/$/,''),managementKey:key,source:source||'unknown'}}
function readCPAConnection(){var manual=loadManualManagementKey();if(manual)return{apiBase:'',managementKey:manual,source:'panel'};var input=byId('mgmtKey');if(input&&String(input.value||'').trim())return{apiBase:'',managementKey:String(input.value).trim(),source:'input'};var raw=storageGet('cli-proxy-auth');var state=decodeCliProxyAuth(raw);var conn=connectionFromState(state,'cli-proxy-auth');if(conn)return conn;try{var legacyKey=storageGet('managementKey');if(legacyKey){var apiBase=storageGet('apiBase')||storageGet('apiUrl')||'';return{apiBase:String(apiBase||'').replace(/\/$/,''),managementKey:String(legacyKey).trim(),source:'legacy'}}}catch(e){}return null}
function updateAuthBadge(){var el=byId('authBadge');if(!el)return;var conn=readCPAConnection();if(conn&&conn.managementKey){el.textContent='管理授权：已就绪（'+conn.source+'）';el.className='auth-badge ok'}else{el.textContent='管理授权：未就绪 — 删除/检查不可用';el.className='auth-badge warn'}}
function managementFetch(path,options){var conn=readCPAConnection();if(!conn||!conn.managementKey)throw new Error('当前没有可用的管理密钥。请在设置里填写 CPA 管理密钥并保存，或在管理中心勾选“记住密码”后重新登录。');options=options||{};options.headers=options.headers||{};options.headers.Authorization='Bearer '+conn.managementKey;if(!options.headers.accept)options.headers.accept='application/json';var base=(conn.apiBase||window.location.origin).replace(/\/$/,'');return fetch(base+'/v0/management'+path,options)}
function setFeedback(msg,type){var el=byId('feedback');el.className='feedback '+(type||'');el.textContent=msg}
function setBusy(flag){busy=!!flag;document.body.classList.toggle('busy',busy);updateToolbarState();renderTable()}
function parseJsonText(text,endpoint){try{return text?JSON.parse(text):{}}catch(e){var low=String(text||'').toLowerCase();if(low.indexOf('<!doctype')>=0||low.indexOf('<html')>=0)throw new Error('操作端点 '+endpoint+' 未启用：当前后端返回了面板页面，请升级插件后端或注册该管理路由。');throw new Error('操作端点 '+endpoint+' 返回非 JSON：'+String(text||'').slice(0,90))}}
async function managementPluginPost(path,payload){var resp=await managementFetch('/plugins/grok-panel/'+String(path).replace(/^\/+/,''),{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify(payload||{})});var text=await resp.text();var data=parseJsonText(text,path);if(!resp.ok)throw new Error('HTTP '+resp.status+'：'+messageFromData(data,text));return data||{}}
async function runPluginChecks(emails){var indices=[];emails.forEach(function(email){var x=accountByEmail(email);var idx=x&&String(x.auth_index||x.authIndex||'').trim();if(idx)indices.push(idx)});var records=[];for(var i=0;i<indices.length;i++){var data=await managementPluginPost('checks',{auth_index:indices[i]});if(Array.isArray(data.records))records=records.concat(data.records)}return{records:records}}
async function deleteAuthNames(names){names=unique(names).map(normalizeAuthFileName).filter(Boolean);if(!names.length)throw new Error('没有有效的 auth 文件名');var lastErr=null;var attempts=[{body:{names:names}},{body:{name:names[0]},onlySingle:true},{query:names}];for(var a=0;a<attempts.length;a++){var attempt=attempts[a];if(attempt.onlySingle&&names.length!==1)continue;try{var url='/auth-files';if(attempt.query){url+='?'+attempt.query.map(function(n){return 'name='+encodeURIComponent(n)}).join('&')}var opts={method:'DELETE',headers:{accept:'application/json'}};if(attempt.body){opts.headers['content-type']='application/json';opts.body=JSON.stringify(attempt.body)}var resp=await managementFetch(url,opts);var text=await resp.text();var data=text?parseJsonText(text,'auth-files'):{};if(resp.ok||resp.status===207)return data||{status:'ok'};lastErr=new Error('HTTP '+resp.status+'：'+messageFromData(data,text));if(resp.status===401||resp.status===403)throw lastErr}catch(e){lastErr=e;if(String(e.message||'').indexOf('401')>=0||String(e.message||'').indexOf('403')>=0)throw e}}throw lastErr||new Error('删除失败')}
function messageFromData(data,text){if(data&&data.error&&data.error.message)return data.error.message;if(data&&data.message)return data.message;if(data&&data.error)return String(data.error);return String(text||'操作失败').slice(0,120)}
async function fetchData(showFeedback,skipAuto){if(showFeedback)setFeedback('正在刷新数据...','');byId('statusDot').classList.remove('err');try{var resp=await fetch(fixedApiUrl('data'),{credentials:'same-origin',headers:{'accept':'application/json'}});var text=await resp.text();if(resp.status===401||resp.status===403)throw new Error('数据接口授权不可用：请检查 CPA 管理会话。');if(!resp.ok)throw new Error('HTTP '+resp.status);var data=parseJsonText(text,'data');if(data&&data.ok===true&&data.result!==undefined)data=data.result;lastData=normalizeData(data);allData=lastData.files||[];allData.forEach(function(x,i){x._rowKey=makeRowKey(x,i)});pruneSelection();renderStats(lastData);renderChart(lastData.recent_buckets||[]);renderTable();byId('lastUpdate').textContent='更新于 '+fmtTime();if(showFeedback)setFeedback('数据已刷新。','ok');if(settings.autoCheck&&!skipAuto)maybeAutoCheck()}catch(e){byId('statusDot').classList.add('err');setFeedback('连接失败：'+e.message,'error')}}
function normalizeData(data){data=data||{};if(!Array.isArray(data.files))data.files=[];return data}
function makeRowKey(x,i){var email=getEmail(x);return email?email.toLowerCase():'row-'+i}
function getEmail(x){return String((x&&(x.email||x.account||''))||'').trim()||String((x&&(x.name||x.id))||'').replace(/\.json$/i,'').trim()}
function normalizeAuthFileName(name){name=String(name||'').trim();if(!name)return'';name=name.split(/[\\/]/).pop();if(!/\.json$/i.test(name))name+='.json';return name}
function getAuthName(x){if(!x)return'';var raw=x.name||x.id||x.file||x.filename||'';raw=String(raw||'').trim();if(raw&&!/\.json$/i.test(raw)&&raw.indexOf('@')<0&&raw.length<=32){/* likely auth_index, skip */}if(raw&&/\.json$/i.test(raw))return normalizeAuthFileName(raw);if(raw&&raw.indexOf('@')>=0)return normalizeAuthFileName(raw);var email=getEmail(x);if(email&&email.indexOf('@')>=0)return normalizeAuthFileName(email);return normalizeAuthFileName(raw)}
function getStatus(x){return String((x&&x.status)||'').trim()}
function tokenLimit(){return Math.max(1,parseInt(settings.tokenLimit,10)||2000000)}
function tokensPerReq(){return Math.max(1,parseInt(settings.tokensPerReq,10)||5000)}
function failThreshold(){return Math.max(1,parseInt(settings.threshold,10)||3)}
function usagePct(x){var tl=tokenLimit();var et=(Number(x&&x.success)||0)*tokensPerReq();return tl>0?Math.max(0,Math.min(999,et/tl*100)):0}
function isMobileView(){return window.matchMedia&&window.matchMedia('(orientation:portrait),(hover:none) and (pointer:coarse)').matches}
function meterWidth(pct,kind){var span=kind==='row'?(isMobileView()?24:11):(isMobileView()?80:18);var clamped=Math.max(0,Math.min(100,Number(pct)||0));return(clamped*span/100).toFixed(2)+'vw'}
function classifyType(x){var key=makeRowKey(x,0);var raw=x.tier||typeOverrides[key]||x.account_type||x.accountType||x.account_kind||x.accountKind||x.plan||x.type||x.label||'';raw=String(raw||'').trim();var low=raw.toLowerCase();if(!raw||low==='unknown'||low==='unk'||low==='oauth'||low==='none'||low==='')return{key:'free',label:'Free'};if(low.indexOf('heavy')>=0||low.indexOf('bulk')>=0||low.indexOf('team')>=0)return{key:'heavy',label:'Heavy'};if(low.indexOf('super')>=0||low.indexOf('premium')>=0||low.indexOf('paid')>=0||low.indexOf('pro')>=0||low.indexOf('max')>=0)return{key:'super',label:'Super'};if(low.indexOf('free')>=0||low.indexOf('basic')>=0||low.indexOf('standard')>=0||low.indexOf('normal')>=0)return{key:'free',label:'Free'};return{key:'free',label:'Free'}}
function isHeavyAccount(x){return classifyType(x).key==='heavy'}
function statusKey(x){var s=getStatus(x).toLowerCase();if(x&&x.disabled)return'disabled';if(!s)return'unknown';if(s.indexOf('disable')>=0||s.indexOf('off')>=0)return'disabled';if(s.indexOf('active')>=0||s.indexOf('ok')>=0||s.indexOf('ready')>=0||s==='available'||s.indexOf('healthy')>=0)return'active';if(s.indexOf('error')>=0||s.indexOf('fail')>=0||s.indexOf('invalid')>=0)return'error';if(s.indexOf('unavail')>=0||s.indexOf('cooling')>=0||s.indexOf('cooldown')>=0||s.indexOf('retry')>=0||s.indexOf('rate')>=0||s.indexOf('quota')>=0)return'unavailable';return'other'}
function mapHealth(raw){var low=String(raw||'').toLowerCase();if(!low)return null;if(low.indexOf('disabled')>=0||low.indexOf('off')>=0)return{key:'disabled',label:'禁用',detail:raw};if(low.indexOf('unavail')>=0||low.indexOf('cooling')>=0||low.indexOf('cooldown')>=0||low.indexOf('retry')>=0||low.indexOf('rate')>=0||low.indexOf('quota')>=0||low==='error')return{key:'unavailable',label:'异常',detail:raw};if(low.indexOf('invalid')>=0||low.indexOf('expired')>=0||low.indexOf('revoked')>=0||low.indexOf('dead')>=0)return{key:'invalid',label:'无效',detail:raw};if(low.indexOf('warn')>=0||low.indexOf('limited')>=0||low.indexOf('fail')>=0)return{key:'warn',label:'警告',detail:raw};if(low.indexOf('healthy')>=0||low.indexOf('active')>=0||low.indexOf('ok')>=0||low.indexOf('valid')>=0)return{key:'healthy',label:'健康',detail:raw};if(low.indexOf('unknown')>=0)return{key:'unknown',label:'未知',detail:raw};return null}
function deriveHealth(x){var key=makeRowKey(x,0);if(healthOverrides[key])return healthOverrides[key];var mapped=mapHealth(x.health||x.account_health||x.accountHealth||x.health_status||x.healthStatus);if(mapped)return mapped;if(x&&x.disabled)return{key:'disabled',label:'禁用',detail:'CPA 已禁用'};if(x&&x.unavailable)return{key:'unavailable',label:'异常',detail:getStatus(x)||'unavailable'};var sk=statusKey(x);if(sk==='active')return{key:'healthy',label:'健康',detail:getStatus(x)||'CPA active'};if(sk==='error')return{key:'unavailable',label:'异常',detail:getStatus(x)||'error'};if(sk==='unavailable')return{key:'unavailable',label:'异常',detail:getStatus(x)||'unavailable'};if(sk==='disabled')return{key:'disabled',label:'禁用',detail:getStatus(x)||'disabled'};return{key:'unknown',label:'未检查',detail:'需要执行健康检查'}}
function statusLabel(sk){return({active:'活跃',disabled:'禁用',error:'错误',unavailable:'异常',other:'其他',unknown:'未知'})[sk]||sk}
function protectReason(x){var t=classifyType(x);if(settings.protectSuper&&t.key==='super')return'保护 super';if(settings.protectHeavy&&isHeavyAccount(x))return'保护 heavy';if(settings.protectUnknown&&t.key==='unknown')return'保护 unknown';return''}
function isInvalidCandidate(x){var h=deriveHealth(x).key;return h==='invalid'||h==='disabled'}
function isDeleteEligible(x){return!!(x&&x.delete_eligible===true&&getEmail(x)&&!protectReason(x))}
function cleanupCandidates(list){return(list||getFilteredData()).filter(isDeleteEligible)}
function setMeter(id,pct){var el=byId(id);if(!el)return;el.style.width=meterWidth(pct,'stat');el.className='bar-fill'+(pct>=80?' danger':pct>=50?' warn':'')}
function renderStats(d){d=d||{};var total=d.total_files||allData.length||0;var active=d.active_files||0;var dis=d.disabled_files||0;var ts=d.total_success||0;var tf=d.total_failed||0;var errorN=allData.filter(function(x){return statusKey(x)==='error'||statusKey(x)==='unavailable'||!!x.unavailable}).length;if(!d.total_files&&allData.length){active=allData.filter(function(x){return statusKey(x)==='active'}).length;dis=allData.filter(function(x){return statusKey(x)==='disabled'}).length;ts=sum(allData,'success');tf=sum(allData,'failed')}var et=ts*tokensPerReq();var cap=total*tokenLimit();var pct=cap>0?et/cap*100:0;var activePct=total>0?active/total*100:0;var reqTotal=ts+tf;var reqPct=reqTotal>0?ts/reqTotal*100:0;var typeCounts={free:0,super:0,heavy:0,unknown:0};var healthCounts={healthy:0,warn:0,invalid:0,disabled:0,unavailable:0,unknown:0};allData.forEach(function(x){var t=classifyType(x).key;typeCounts[t]=(typeCounts[t]||0)+1;var h=deriveHealth(x).key;healthCounts[h]=(healthCounts[h]||0)+1});var healthy=healthCounts.healthy||0;var healthPct=total>0?healthy/total*100:0;var invalidAll=allData.filter(isInvalidCandidate);var cleanAll=cleanupCandidates(allData);byId('statTotal').textContent=fmt(total);byId('statTotalSub').textContent='共 '+fmt(total)+' 个文件，'+fmt(dis)+' 个已禁用';byId('statActive').textContent=activePct.toFixed(1)+'%';byId('statActiveSub').textContent='活跃 '+fmt(active)+' / '+fmt(total)+'，异常/错误 '+fmt(errorN);setMeter('activeBar',activePct);byId('statRequests').textContent=reqPct.toFixed(1)+'%';byId('statRequestsSub').textContent='成功 '+fmt(ts)+' / '+fmt(reqTotal)+'，失败 '+fmt(tf);setMeter('requestBar',reqPct);byId('statUsage').textContent=pct.toFixed(2)+'%';byId('statUsageSub').textContent='已用 '+fmt(et)+' / 容量 '+fmt(cap);setMeter('usageBar',pct);byId('statTokens').textContent=fmt(et);byId('statTokensSub').textContent='按 '+fmt(tokensPerReq())+' token/请求估算';byId('statCapacity').textContent=fmt(cap);byId('statCapacitySub').textContent='上限 '+fmt(tokenLimit())+' token/账号 × '+fmt(total);byId('statTypes').textContent='Free '+fmt(typeCounts.free)+' · Super '+fmt(typeCounts.super)+' · Heavy '+fmt(typeCounts.heavy);byId('statTypesSub').textContent='未知类型 '+fmt(typeCounts.unknown)+' 个';byId('statHealth').textContent=healthPct.toFixed(1)+'%';byId('statHealthSub').textContent='健康 '+fmt(healthy)+' / '+fmt(total)+'，异常 '+fmt(healthCounts.unavailable)+'，警告 '+fmt(healthCounts.warn)+'，无效 '+fmt(healthCounts.invalid);setMeter('healthBar',healthPct);byId('statInvalid').textContent=fmt(invalidAll.length);byId('statInvalidSub').textContent=fmt(cleanAll.length)+' 已满足连续 '+failThreshold()+' 次 401/403，可安全清理；其余 '+fmt(Math.max(0,invalidAll.length-cleanAll.length))+' 个仅标记异常/无效';byId('statSelected').textContent=fmt(getSelectedEmails().length);byId('statSelectedSub').textContent=fmt(getFilteredData().length)+' 个当前可见'}
function sum(arr,key){return arr.reduce(function(n,x){return n+(Number(x&&x[key])||0)},0)}
function renderChart(buckets){var c=byId('chartRow');c.innerHTML='';if(!buckets||!buckets.length){c.innerHTML='<div class="chart-empty">暂无趋势数据</div>';return}buckets=buckets.slice().sort(function(a,b){return a.time<b.time?-1:1});var max=1;buckets.forEach(function(b){max=Math.max(max,(Number(b.success)||0)+(Number(b.failed)||0))});buckets.forEach(function(b){var total=(Number(b.success)||0)+(Number(b.failed)||0);var h=Math.max(0.35,total/max*14);var bar=document.createElement('div');bar.className='chart-bar'+((Number(b.failed)||0)>(Number(b.success)||0)?' fail':'');bar.style.height=h.toFixed(2)+'vh';bar.setAttribute('data-tip',String(b.time||'')+' | 成功:'+fmt(b.success)+' 失败:'+fmt(b.failed));c.appendChild(bar)})}
function pageSize(){var n=parseInt(settings.pageSize,10);if([20,50,100,200].indexOf(n)<0)n=50;return n}
function getPageMeta(list){list=list||[];var size=pageSize();var total=list.length;var totalPages=Math.max(1,Math.ceil(total/size)||1);if(currentPage>totalPages)currentPage=totalPages;if(currentPage<1)currentPage=1;var start=(currentPage-1)*size;var end=Math.min(total,start+size);return{list:list,size:size,total:total,totalPages:totalPages,page:currentPage,start:start,end:end,rows:list.slice(start,end)}}
function renderPager(meta){var info='第 <b>'+meta.page+'</b> / '+meta.totalPages+' 页 · 每页 '+meta.size;var range=meta.total?('显示 '+(meta.start+1)+'-'+meta.end+' / '+meta.total):'显示 0-0 / 0';['pageInfo','pageInfo2'].forEach(function(id){if(byId(id))byId(id).innerHTML=info});if(byId('pageRange'))byId('pageRange').textContent=range;var atFirst=meta.page<=1;var atLast=meta.page>=meta.totalPages;[['pageFirstBtn',atFirst],['pagePrevBtn',atFirst],['pagePrevBtn2',atFirst],['pageNextBtn',atLast],['pageNextBtn2',atLast],['pageLastBtn',atLast]].forEach(function(pair){var el=byId(pair[0]);if(el)el.disabled=!!busy||!!pair[1]})}
function getFilteredData(){var s=byId('searchBox')?byId('searchBox').value.toLowerCase().trim():'';var sf=byId('statusFilter')?byId('statusFilter').value:'all';var tf=byId('typeFilter')?byId('typeFilter').value:'all';var hf=byId('healthFilter')?byId('healthFilter').value:'all';var uf=byId('usageFilter')?byId('usageFilter').value:'all';var arr=allData.filter(function(x){var email=getEmail(x).toLowerCase();var st=getStatus(x).toLowerCase();var typ=classifyType(x);var health=deriveHealth(x);var pct=usagePct(x);if(s&&email.indexOf(s)<0&&st.indexOf(s)<0&&typ.label.toLowerCase().indexOf(s)<0&&health.label.indexOf(s)<0)return false;if(sf!=='all'&&statusKey(x)!==sf)return false;if(tf!=='all'&&typ.key!==tf)return false;if(hf!=='all'&&health.key!==hf)return false;if(uf==='unused'&&(Number(x.success)||0)!==0)return false;if(uf==='low'&&pct>=50)return false;if(uf==='warn'&&(pct<50||pct>=80))return false;if(uf==='high'&&pct<80)return false;return true});var sort=byId('sortFilter')?byId('sortFilter').value:'success_desc';arr.sort(function(a,b){if(sort==='failed_desc')return(Number(b.failed)||0)-(Number(a.failed)||0);if(sort==='usage_desc')return usagePct(b)-usagePct(a);if(sort==='health_asc')return healthRank(a)-healthRank(b);if(sort==='type_asc')return classifyType(a).key<classifyType(b).key?-1:1;if(sort==='email_asc')return getEmail(a)<getEmail(b)?-1:1;return(Number(b.success)||0)-(Number(a.success)||0)});return arr}
function healthRank(x){var order={invalid:0,disabled:1,unavailable:2,error:2,warn:3,unknown:4,healthy:5};return order[deriveHealth(x).key]===undefined?9:order[deriveHealth(x).key]}
function renderTable(){var tbody=byId('tableBody');if(!tbody)return;var f=getFilteredData();var meta=getPageMeta(f);renderPager(meta);tbody.innerHTML='';meta.rows.forEach(function(x,i){var email=getEmail(x);var absIndex=meta.start+i;var key=makeRowKey(x,absIndex);var su=Number(x.success)||0;var fa=Number(x.failed)||0;var totalReq=su+fa;var successPct=totalReq>0?su/totalReq*100:0;var et=su*tokensPerReq();var tl=tokenLimit();var pct=usagePct(x);var bc=pct>80?'danger':pct>50?'warn':'';var typ=classifyType(x);var h=deriveHealth(x);var prot=protectReason(x);var armed=isArmed('delete:'+key);var disabled=busy||!email||!!prot;var tr=document.createElement('tr');tr.className='row-'+h.key;var protectedText=prot?'<div class="cell-sub">⚠ '+esc(prot)+'</div>':'';var healthDetail=h.detail?'<div class="cell-sub">'+esc(h.detail)+'</div>':'';var statusSub=getStatus(x)||'-';var accountHtml='<div class="email-main">'+esc(email||'?')+'</div><div class="email-meta"><span class="tag '+esc(statusKey(x))+'">'+esc(statusLabel(statusKey(x)))+'</span><span class="cell-sub">'+esc(statusSub)+'</span></div>'+protectedText;var healthHtml='<span class="health-indicator"><span class="health-dot '+esc(h.key)+'"></span><span class="tag '+esc(h.key)+'">'+esc(h.label)+'</span></span>'+healthDetail;var usageHtml='<div class="usage-cell"><div class="metric-compact"><span>成功 <b>'+fmt(su)+'</b></span><span>失败 <b class="'+(fa>0?'red-text':'')+'">'+fmt(fa)+'</b></span></div><div class="metric-compact"><span>成功率 <b>'+successPct.toFixed(0)+'%</b></span><span><b>'+fmt(et)+'</b>/'+fmt(tl)+'</span></div><div class="bar-container small"><div class="bar-fill '+bc+'" style="width:'+meterWidth(pct,'row')+'"></div></div></div>';tr.innerHTML='<td class="col-check"><input type="checkbox" class="row-select" data-key="'+esc(key)+'" data-email="'+esc(email)+'" '+(selected[key]?'checked':'')+' '+(!email||busy?'disabled':'')+'></td><td class="col-idx num-cell">'+(absIndex+1)+'</td><td class="email-cell">'+accountHtml+'</td><td class="col-type"><span class="tag '+esc(typ.key)+'">'+esc(typ.label)+'</span></td><td class="col-health">'+healthHtml+'</td><td class="col-usage">'+usageHtml+'</td><td class="col-actions"><div class="actions-cell"><button data-act="check" data-key="'+esc(key)+'" '+(busy||!email?'disabled':'')+'>检查</button><button data-act="verify-tier" data-key="'+esc(key)+'" '+(busy||!x.auth_index?'disabled':'')+'>核实</button><button data-act="delete" data-key="'+esc(key)+'" '+(disabled?'disabled':'')+' class="'+(armed?'armed':'')+'">'+(armed?'确认':'删除')+'</button></div></td>';tbody.appendChild(tr)});if(!meta.rows.length)tbody.innerHTML='<tr><td colspan="7" class="muted" style="padding:3vh;text-align:center">没有匹配的账号</td></tr>';renderSummary(f);updateToolbarState(meta.rows)}

function renderStatsOnlySelection(){byId('statSelected').textContent=getSelectedEmails().length;byId('statSelectedSub').textContent=getFilteredData().length+' 个当前可见'}
function renderSummary(f){var used=f.filter(function(x){return(Number(x.success)||0)>0}).length;var avail=f.filter(function(x){return(Number(x.success)||0)===0&&!x.disabled}).length;var ts=sum(f,'success');var tf=sum(f,'failed');var invalid=f.filter(isInvalidCandidate).length;var clean=cleanupCandidates(f).length;byId('sumCount').textContent=f.length;byId('sumUsed').textContent=used;byId('sumAvail').textContent=avail;byId('sumRate').textContent=(ts+tf>0?(ts/(ts+tf)*100).toFixed(1):0)+'%';byId('sumInvalid').textContent=invalid;byId('sumCleanup').textContent=clean;var meta=getPageMeta(f);byId('filterInfo').textContent=f.length+' / '+allData.length+' 过滤，本页 '+(meta.total?(meta.start+1)+'-'+meta.end:'0')+'，'+clean+' 已满足安全清理条件'}
function renderAll(){if(lastData){renderStats(lastData);renderChart(lastData.recent_buckets||[])}renderTable()}
async function verifyTiers(accounts,label){accounts=(accounts||[]).filter(function(x){return x&&x.auth_index});if(!accounts.length){setFeedback('没有具备 auth_index 的账号可核实。','warn');return}if(busy)return;setBusy(true);var ok=0,unknown=0,failed=[];try{for(var i=0;i<accounts.length;i++){var x=accounts[i];setFeedback(label+'：'+(i+1)+' / '+accounts.length+' · '+getEmail(x),'');try{var data=await managementPluginPost('verify-tier',{auth_index:x.auth_index});var results=applyCheckResult(data);var c=results[0]&&results[0].classification||{};if(c.tier&&c.tier!=='unknown')ok++;else{unknown++;if(c.detail)failed.push(getEmail(x)+'：'+c.detail)}}catch(e){unknown++;failed.push(getEmail(x)+'：'+e.message)}renderAll()}await fetchData(false,true);var msg=label+'完成：明确识别 '+ok+' 个，仍未知 '+unknown+' 个。';if(failed.length)msg+=' 首个原因：'+failed[0];setFeedback(msg,unknown?'warn':'ok')}finally{setBusy(false)}}
function verifyVisibleTiers(){verifyTiers(getPageMeta(getFilteredData()).rows,'核实本页套餐')}
function verifySelectedTiers(){var emails=getSelectedEmails();var rows=emails.map(accountByEmail).filter(Boolean);if(!rows.length){setFeedback('请先选择账号。','warn');return}verifyTiers(rows,'核实选中套餐')}
async function verifyTierOne(key){var x=accountByKey(key);if(!x||!x.auth_index){setFeedback('该账号没有 auth_index，无法核实套餐。','warn');return}setBusy(true);try{var data=await managementPluginPost('verify-tier',{auth_index:x.auth_index});var results=applyCheckResult(data);var r=results[0]||{};var c=r.classification||{};setFeedback('套餐核实完成：'+String(c.tier||'unknown')+' · '+String(c.detail||c.source||'已更新'),'ok');renderAll();fetchData(false,true)}catch(e){setFeedback('套餐核实失败：'+e.message,'error')}finally{setBusy(false)}}
function handleTableClick(evt){var btn=evt.target.closest('button[data-act]');if(!btn)return;var key=btn.getAttribute('data-key');if(btn.getAttribute('data-act')==='check')checkOne(key);if(btn.getAttribute('data-act')==='verify-tier')verifyTierOne(key);if(btn.getAttribute('data-act')==='delete')requestSingleDelete(key)}
function handleTableChange(evt){if(!evt.target.classList.contains('row-select'))return;var key=evt.target.getAttribute('data-key');var email=evt.target.getAttribute('data-email');if(evt.target.checked)selected[key]=email;else delete selected[key];renderAll()}
function toggleVisibleSelection(){var on=byId('selectVisible').checked;var meta=getPageMeta(getFilteredData());meta.rows.forEach(function(x,i){var email=getEmail(x);if(!email)return;var key=makeRowKey(x,meta.start+i);if(on)selected[key]=email;else delete selected[key]});renderAll()}
function pruneSelection(){var valid={};allData.forEach(function(x,i){valid[makeRowKey(x,i)]=true});Object.keys(selected).forEach(function(k){if(!valid[k])delete selected[k]})}
function getSelectedEmails(){var out=[];Object.keys(selected).forEach(function(k){if(selected[k])out.push(selected[k])});return unique(out)}
function unique(arr){var seen={};var out=[];(arr||[]).forEach(function(v){v=String(v||'').trim();var k=v.toLowerCase();if(v&&!seen[k]){seen[k]=true;out.push(v)}});return out}
function updateToolbarState(f){if(!byId('selectionInfo'))return;var filtered=getFilteredData();var meta=getPageMeta(filtered);f=f||meta.rows;var count=getSelectedEmails().length;var visibleWithEmail=f.filter(function(x){return!!getEmail(x)});var selectedVisible=0;f.forEach(function(x,i){var key=makeRowKey(x,meta.start+i);if(getEmail(x)&&selected[key])selectedVisible++});var cb=byId('selectVisible');cb.checked=visibleWithEmail.length>0&&selectedVisible===visibleWithEmail.length;cb.indeterminate=selectedVisible>0&&selectedVisible<visibleWithEmail.length;byId('selectionInfo').textContent=count?('已选择 '+count+' 个'):'未选择';byId('checkSelectedBtn').disabled=busy||count===0;byId('verifySelectedBtn').disabled=busy||count===0;byId('verifyVisibleTopBtn').disabled=busy||meta.rows.filter(function(x){return!!x.auth_index}).length===0;byId('batchDeleteBtn').disabled=busy||count===0;byId('batchDeleteBtn').textContent=isArmed('delete:batch')?'确认删除选中':'删除选中';if(isArmed('delete:batch'))byId('batchDeleteBtn').classList.add('armed');else byId('batchDeleteBtn').classList.remove('armed');var clean=cleanupCandidates(filtered).length;byId('cleanupInvalidBtn').disabled=busy||clean===0;byId('cleanupInvalidBtn').textContent=isArmed('cleanup:visible')?'确认安全清理':'安全清理 '+clean;if(isArmed('cleanup:visible'))byId('cleanupInvalidBtn').classList.add('armed');else byId('cleanupInvalidBtn').classList.remove('armed');byId('checkVisibleTopBtn').disabled=busy||meta.rows.filter(function(x){return!!getEmail(x)}).length===0;byId('refreshBtn').disabled=busy}
function isArmed(key){return(confirmUntil[key]||0)>Date.now()}
function armOrProceed(key,msg){if(isArmed(key)){delete confirmUntil[key];return true}confirmUntil[key]=Date.now()+6500;setFeedback(msg,'warn');window.setTimeout(function(){if(confirmUntil[key]&&confirmUntil[key]<Date.now()){delete confirmUntil[key];renderAll()}},6800);renderAll();return false}
function accountByKey(key){for(var i=0;i<allData.length;i++){if(makeRowKey(allData[i],i)===key)return allData[i]}return null}
function accountByEmail(email){var target=String(email||'').toLowerCase();for(var i=0;i<allData.length;i++){if(getEmail(allData[i]).toLowerCase()===target)return allData[i]}return null}
function checkOne(key){var x=accountByKey(key);var email=x?getEmail(x):'';if(!email){setFeedback('无法检查：账号缺少邮箱或 ID。','error');return}runHealthCheck([email],'手动检查 1 个',false)}
function manualCheckVisible(){var emails=unique(getPageMeta(getFilteredData()).rows.map(getEmail));if(!emails.length){setFeedback('没有可检查的本页账号。','warn');return}runHealthCheck(emails,'手动检查本页 '+emails.length+' 个',false)}
function manualCheckSelected(){var emails=getSelectedEmails();if(!emails.length){setFeedback('请先选择账号。','warn');return}runHealthCheck(emails,'手动检查选中 '+emails.length+' 个',false)}
async function runHealthCheck(emails,label,autoMode){emails=unique(emails);if(!emails.length)return;if(busy&&!autoMode)return;setBusy(true);try{var data=await runPluginChecks(emails);var results=applyCheckResult(data);setFeedback(label+' 完成：收到 '+results.length+' 条检查结果。','ok');if(settings.autoDelete)await runAutoCleanup();renderAll();fetchData(false,true)}catch(e){setFeedback(label+' 失败：'+e.message,'error')}finally{setBusy(false)}}
function applyCheckResult(data){var results=[];if(data){if(Array.isArray(data.records))results=data.records;else if(Array.isArray(data.results))results=data.results;else if(Array.isArray(data.files))results=data.files;else if(Array.isArray(data.accounts))results=data.accounts;else if(Array.isArray(data.checked))results=data.checked;else if(typeof data==='object'){Object.keys(data).forEach(function(k){if(data[k]&&typeof data[k]==='object'){var v=data[k];if(!v.email)v.email=k;results.push(v)}})}}results.forEach(function(r){var email=String(r.email||r.account||r.id||'').trim();if(!email)return;var key=email.toLowerCase();var mapped=mapHealth(r.health||r.status||r.result||'')||{key:'unknown',label:'未知',detail:r.reason||r.detail||r.message||'checked'};mapped.detail=String(r.reason||r.detail||r.message||mapped.detail||'checked');healthOverrides[key]=mapped;var classification=r.classification||{};var tier=classification.tier;if(tier||r.account_type||r.accountType||r.type)typeOverrides[key]=tier||r.account_type||r.accountType||r.type;var x=accountByEmail(email);if(x){x.delete_eligible=r.delete_eligible===true;x.invalid_streak=Number(r.invalid_streak)||0;x.health=r.health||x.health;x.tier=tier||x.tier;x.tier_source=classification.source||x.tier_source;x.tier_detail=classification.detail||x.tier_detail}});return results}
function requestSingleDelete(key){var x=accountByKey(key);if(!x){setFeedback('无法删除：找不到该账号。','error');return}var email=getEmail(x);var prot=protectReason(x);if(prot){setFeedback('删除已阻止：'+email+' 受 '+prot+' 保护。','warn');return}if(!armOrProceed('delete:'+key,'再次点击该行“确认删除”将删除 '+email+'。'))return;deleteAccounts([email],'删除 '+email)}
function requestBatchDelete(){var emails=getSelectedEmails();if(!emails.length){setFeedback('请先选择账号。','warn');return}if(!armOrProceed('delete:batch','再次点击“确认删除选中”将删除未受保护的选中账号。'))return;deleteAccounts(emails,'批量删除')}
async function deleteAccounts(emails,label){emails=unique(emails);var names=[];var blocked=[];emails.forEach(function(email){var x=accountByEmail(email);var prot=x?protectReason(x):(settings.protectUnknown?'保护 unknown':'');var name=x?getAuthName(x):'';if(prot)blocked.push(email+'('+prot+')');else if(name)names.push(name);else blocked.push(email+'(缺少 auth 文件名)')});names=unique(names.map(normalizeAuthFileName).filter(Boolean));if(!names.length){setFeedback(label+' 已取消：没有未受保护且具备文件名的账号。'+(blocked.length?' 已阻止 '+blocked.length+' 个。':''),'warn');return}if(!readCPAConnection()){setFeedback(label+' 失败：没有管理密钥。请在设置里填写并保存 CPA 管理密钥。','error');updateAuthBadge();return}setBusy(true);try{var data=await deleteAuthNames(names);emails.forEach(function(email){delete selected[email.toLowerCase()]});setFeedback(label+' 完成：已请求删除 '+names.length+' 个 auth 文件。'+mutationSummary(data)+(blocked.length?' 阻止 '+blocked.length+' 个。':''),'ok');await fetchData(false,true)}catch(e){setFeedback(label+' 失败：'+e.message,'error');updateAuthBadge()}finally{setBusy(false)}}
function requestCleanupInvalid(){var candidates=cleanupCandidates(getFilteredData());if(!candidates.length){setFeedback('当前筛选范围可安全清理为 0：只有连续 '+failThreshold()+' 次明确返回 401/403、且未受保护的账号才会进入清理。请先勾选异常账号并点击“检查选中”。','warn');return}if(!armOrProceed('cleanup:visible','再次点击“确认安全清理”将删除当前筛选范围内 '+candidates.length+' 个已满足连续 '+failThreshold()+' 次 401/403 的账号。'))return;cleanupInvalid(candidates,'安全清理')}
async function cleanupInvalid(candidates,label){var emails=unique((candidates||[]).filter(isDeleteEligible).map(getEmail));if(!emails.length){setFeedback('当前没有达到连续 '+failThreshold()+' 次明确 401/403 且未受保护的账号。','warn');return}await deleteAccounts(emails,label)}
function currentProtectSettings(){return{super:!!settings.protectSuper,heavy:!!settings.protectHeavy,unknown:!!settings.protectUnknown}}
function mutationSummary(data){if(!data)return'';var parts=[];if(data.deleted!==undefined)parts.push('已删除 '+data.deleted);if(data.status)parts.push('状态 '+data.status);if(Array.isArray(data.files)&&data.files.length)parts.push('文件 '+data.files.length);if(data.skipped!==undefined)parts.push('跳过 '+data.skipped);if(data.checked!==undefined)parts.push('检查 '+data.checked);if(data.message)parts.push(data.message);return parts.length?'（'+parts.join('，')+'）':''}
function maybeAutoCheck(){if(autoCheckBusy)return;var now=Date.now();if(now-lastAutoCheckAt<300000)return;var emails=unique(allData.filter(function(x){return getEmail(x)&&!x.disabled}).map(getEmail));if(!emails.length)return;lastAutoCheckAt=now;autoCheckBusy=true;runHealthCheck(emails,'自动检查 '+emails.length+' 个',true).finally(function(){autoCheckBusy=false})}
async function runAutoCleanup(){if(!settings.autoDelete)return;var candidates=cleanupCandidates(allData);if(!candidates.length){setFeedback('自动检查完成：没有未受保护的无效账号。','ok');return}await cleanupInvalid(candidates,'自动删除无效')}
function init(){initSettings();bindEvents();var initialView='stats';try{initialView=sessionStorage.getItem('grok-panel-active-view')||'stats'}catch(e){}switchPanel(initialView);updateAuthBadge();fetchData(true,false);refreshHandle=window.setInterval(function(){fetchData(false,false)},30000)}
init();
</script>
</body></html>`
