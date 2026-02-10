package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/PipeOpsHQ/agent-sdk-go/devui"
	"github.com/PipeOpsHQ/agent-sdk-go/flow"
)

const openclawFlowName = "openclaw-ui-example"

func main() {
	addr := flag.String("addr", "0.0.0.0:8091", "UI listen address")
	apiBase := flag.String("api-base", "http://127.0.0.1:7070", "Framework DevUI API base URL")
	apiKey := flag.String("api-key", "", "Optional DevUI API key")
	startAPI := flag.Bool("start-api", false, "Start embedded DevUI API (SDK-style) with openclaw flow")
	apiAddr := flag.String("api-addr", "0.0.0.0:7070", "Embedded DevUI API listen address when --start-api=true")
	flag.Parse()

	if *startAPI {
		registerOpenClawFlow()
		go func() {
			if err := devui.Start(context.Background(), devui.Options{
				Addr:        strings.TrimSpace(*apiAddr),
				DefaultFlow: openclawFlowName,
			}); err != nil {
				log.Fatalf("embedded DevUI API failed: %v", err)
			}
		}()
		*apiBase = "http://" + strings.TrimSpace(*apiAddr)
	}

	html := openclawHTML(strings.TrimRight(*apiBase, "/"), strings.TrimSpace(*apiKey))
	upstream, err := url.Parse(strings.TrimRight(*apiBase, "/"))
	if err != nil {
		log.Fatalf("invalid --api-base URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	})
	mux.Handle("/api/", proxy)

	log.Printf("OpenClaw UI listening on http://%s", *addr)
	log.Printf("Connected API base: %s", *apiBase)
	if *startAPI {
		log.Printf("Embedded DevUI API started on http://%s with default flow %q", *apiAddr, openclawFlowName)
	}
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func registerOpenClawFlow() {
	flow.MustRegister(&flow.Definition{
		Name:        openclawFlowName,
		Description: "OpenClaw-style autonomous profile with kickoff, priorities, and continuous execution updates.",
		Workflow:    "router",
		Tools:       []string{"@all"},
		Skills:      []string{"research-planner", "release-readiness", "oncall-triage", "document-manager"},
		SystemPrompt: `You are OpenClaw Bot, an autonomous execution agent.
- On first engagement in a new thread, ask for your preferred name, top priorities, success criteria, and risk boundaries.
- If this context already exists in the thread, do not re-ask.
- Keep a running status update format: Completed / In progress / Next / Blockers.
- Execute proactively with tools and verify each major step before moving on.
- Keep responses concise and operationally clear.`,
		InputExample: "Take this incident triage backlog and autonomously drive it to resolution with checkpoints.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{"type": "string", "description": "Task, incident, or objective for OpenClaw execution."},
			},
			"required": []string{"input"},
		},
	})
}

func openclawHTML(apiBase, apiKey string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>OpenClaw Demo UI</title>
  <style>
    :root{--bg:#071521;--bg2:#0d2234;--fg:#d9ecff;--muted:#89a8c8;--card:#0f2538;--line:#2a4a68;--accent:#2f9cff;--ok:#39d2a3}
    html,body{height:100%%}
    body{margin:0;background:radial-gradient(circle at 10%% 8%%,#14395a,transparent 40%%),linear-gradient(180deg,var(--bg),#081726);color:var(--fg);font-family:"IBM Plex Sans",ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,Arial;overflow:hidden}
    .app{height:100vh;padding:14px;box-sizing:border-box}
    .shell{max-width:1050px;height:100%%;margin:0 auto;display:grid;grid-template-columns:290px 1fr;gap:12px}
    .panel{background:rgba(15,37,56,.88);border:1px solid var(--line);border-radius:12px;display:flex;flex-direction:column;min-height:0}
    .side{padding:14px;overflow:auto}
    .side h2{margin:0 0 6px;font-size:15px}
    .side p{margin:0 0 10px;color:var(--muted);font-size:12px;line-height:1.45}
    .chip{display:inline-flex;padding:2px 8px;border:1px solid #3f6486;border-radius:999px;font-size:11px;color:#a9c8e8;background:#102a41;margin:0 6px 6px 0}
    .status{font-size:12px;color:#bcd7f3;line-height:1.45;border:1px solid #345673;border-radius:9px;padding:8px;background:#102638}
    .chatHead{padding:12px 14px;border-bottom:1px solid var(--line);display:flex;justify-content:space-between;align-items:center}
    .chatHead h1{margin:0;font-size:15px}
    .meta{font-size:11px;color:var(--muted)}
    .chat{flex:1;overflow:auto;padding:12px;min-height:0;background:linear-gradient(180deg,#0d2031,#0b1b2a)}
    .bubble{border:1px solid #2b4b67;background:#0f2639;border-radius:10px;padding:10px 11px;margin-bottom:10px}
    .bubble.user{background:#123454;border-color:#3a6a90}
    .role{font-size:11px;font-weight:700;color:#8fc0eb;margin-bottom:4px;text-transform:uppercase;letter-spacing:.04em}
    .content{white-space:pre-wrap;line-height:1.48;font-size:14px}
    .meta2{margin-top:6px;font-size:11px;color:#7ea4c9}
    .compose{display:flex;gap:8px;padding:10px;border-top:1px solid var(--line);background:#0c1f2f}
    textarea{flex:1;min-height:46px;max-height:140px;border:1px solid #345776;border-radius:8px;padding:10px;background:#0b1a29;color:var(--fg);resize:vertical}
    button{border:none;border-radius:8px;padding:10px 14px;background:var(--accent);color:#fff;cursor:pointer}
    button:disabled{opacity:.65;cursor:not-allowed}
    .hint{padding:8px 12px;border-bottom:1px solid var(--line);font-size:12px;color:#9bc0e2;background:#0e2233}
    @media (max-width:900px){.shell{grid-template-columns:1fr}.side{display:none}}
  </style>
</head>
<body>
  <div class="app">
    <div class="shell">
      <aside class="panel side">
        <h2>OpenClaw Session</h2>
        <p>Autonomous chat demo using a dedicated SDK flow and streaming updates.</p>
        <div class="chip">Flow: %s</div>
        <div class="chip">Workflow: router</div>
        <div class="chip">Tools: @all</div>
        <div class="chip">Threaded Session</div>
        <div style="height:8px"></div>
        <div class="status" id="statusBox">Waiting for first message.</div>
      </aside>
      <main class="panel">
        <div class="chatHead">
          <h1>OpenClaw Bot (Example)</h1>
          <div class="meta" id="metaLine">ready</div>
        </div>
        <div class="hint">Tip: first message should include priorities so the bot can establish execution mode quickly.</div>
        <div id="chat" class="chat"></div>
        <div class="compose">
          <textarea id="input" placeholder="Give OpenClaw an objective..."></textarea>
          <button id="send">Send</button>
        </div>
      </main>
    </div>
  </div>
<script>
  const API_KEY = %q;
  const FLOW = %q;
  const API_BASE = "/api";
  let sessionId = "";
  const history = [];

  function headers(){const h={"Content-Type":"application/json"}; if(API_KEY) h["X-API-Key"]=API_KEY; return h;}
  function esc(s){return String(s||"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;");}
  function push(role,text,meta=""){
    const chat=document.getElementById("chat");
    const row=document.createElement("div");
    row.className="bubble "+role;
    row.innerHTML='<div class="role">'+(role==='user'?'You':'OpenClaw')+'</div><div class="content"></div>'+(meta?'<div class="meta2">'+meta+'</div>':'');
    row.querySelector('.content').textContent=String(text||"");
    chat.appendChild(row); chat.scrollTop=chat.scrollHeight;
  }

  async function send(){
    const input=document.getElementById("input");
    const btn=document.getElementById("send");
    const text=String(input.value||"").trim();
    if(!text) return;
    input.value="";
    push("user",text);
    history.push({role:"user",content:text});
    btn.disabled=true;
    document.getElementById("statusBox").textContent="Running...";
    document.getElementById("metaLine").textContent="streaming";

    const payload={
      input:text,
      sessionId:sessionId||undefined,
      history,
      flow:FLOW,
      workflow:"router",
      tools:["@all"],
      skills:["research-planner","oncall-triage","release-readiness","document-manager"],
      replyTo:{channel:"devui",destination:"openclaw-ui",threadId:sessionId||"openclaw-demo",metadata:{tab:"openclaw-ui",ui:"example"}}
    };

    try{
      const resp=await fetch(API_BASE+"/v1/playground/stream",{method:"POST",headers:headers(),body:JSON.stringify(payload)});
      if(!resp.ok) throw new Error("HTTP "+resp.status);
      const reader=resp.body.getReader();
      const decoder=new TextDecoder();
      let buffer="",evt="",acc=[]; let out=""; let complete=null;
      const streamBubble=document.createElement("div");
      streamBubble.className="bubble assistant";
      streamBubble.innerHTML='<div class="role">OpenClaw</div><div class="content"></div><div class="meta2" id="liveMeta">thinking...</div>';
      document.getElementById("chat").appendChild(streamBubble);
      const c=streamBubble.querySelector('.content');
      const m=streamBubble.querySelector('#liveMeta');
      while(true){
        const {done,value}=await reader.read(); if(done) break;
        buffer += decoder.decode(value,{stream:true});
        const lines=buffer.split("\n"); buffer=lines.pop()||"";
        for(const line of lines){
          if(line.startsWith("event: ")) evt=line.slice(7).trim();
          else if(line.startsWith("data: ")) acc.push(line.slice(6));
          else if(line===""){
            if(evt && acc.length){
              try{
                const data=JSON.parse(acc.join("\n"));
                if(evt==="delta"){out += String(data.text||""); c.textContent=out;}
                else if(evt==="progress"){m.textContent=[data.kind,data.status,data.name||data.toolName||""].filter(Boolean).join(" • ");}
                else if(evt==="complete"){complete=data;}
              }catch(_){ }
            }
            evt=""; acc=[];
          }
        }
      }
      streamBubble.remove();
      if(!complete) throw new Error("stream ended without completion");
      if(complete.sessionId) sessionId = complete.sessionId;
      const meta=[complete.provider?"provider="+complete.provider:"",complete.runId?"run="+complete.runId:"",complete.sessionId?"session="+complete.sessionId:""]
        .filter(Boolean).join(" • ");
      const textOut=complete.output||complete.error||"(empty response)";
      push("assistant",textOut,meta);
      history.push({role:"assistant",content:textOut});
      if(history.length>30) history.splice(0, history.length-30);
      document.getElementById("statusBox").textContent="Completed last run.";
      document.getElementById("metaLine").textContent=meta||"done";
    }catch(e){
      push("assistant","Request failed: "+(e.message||e));
      document.getElementById("statusBox").textContent="Failed.";
      document.getElementById("metaLine").textContent="error";
    }finally{btn.disabled=false;}
  }

  document.getElementById("send").addEventListener("click",send);
  document.getElementById("input").addEventListener("keydown",(e)=>{if(e.key==="Enter"&&!e.shiftKey){e.preventDefault();send();}});
  push("assistant","OpenClaw online. First, tell me what to call myself, your top priorities, and risk boundaries for this thread.");
</script>
</body>
</html>`, openclawFlowName, apiKey, openclawFlowName)
}
