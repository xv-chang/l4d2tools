import{r as e,o as t,c as s,a as r,w as o,b as n,d as a,e as l,f as c}from"./vendor.b62f66d0.js";!function(e=".",t="__import__"){try{self[t]=new Function("u","return import(u)")}catch(s){const r=new URL(e,location),o=e=>{URL.revokeObjectURL(e.src),e.remove()};self[t]=e=>new Promise(((s,n)=>{const a=new URL(e,r);if(self[t].moduleMap[a])return s(self[t].moduleMap[a]);const l=new Blob([`import * as m from '${a}';`,`${t}.moduleMap['${a}']=m;`],{type:"text/javascript"}),c=Object.assign(document.createElement("script"),{type:"module",src:URL.createObjectURL(l),onerror(){n(new Error(`Failed to import: ${e}`)),o(c)},onload(){s(self[t].moduleMap[a]),o(c)}});document.head.appendChild(c)})),self[t].moduleMap={}}}("/assets/");const i={class:"paper container"},p=r("div",{class:"row flex-center"},[r("div",{class:"text-center"},[r("h2",null,"l4d2 Tools")])],-1),d={class:"row flex-center"},m=n("地图下载"),u=n("RCON命令"),f=n("设置"),_={expose:[],setup:n=>(n,a)=>{const l=e("router-link"),c=e("router-view");return t(),s("div",i,[p,r("div",d,[r(l,{class:"paper-btn margin",to:"/maps"},{default:o((()=>[m])),_:1}),r(l,{class:"paper-btn margin",to:"/mods"},{default:o((()=>[u])),_:1}),r(l,{class:"paper-btn margin",to:"/settings"},{default:o((()=>[f])),_:1})]),r(c)])}};let h;const v={},E=function(e,t){if(!t)return e();if(void 0===h){const e=document.createElement("link").relList;h=e&&e.supports&&e.supports("modulepreload")?"modulepreload":"preload"}return Promise.all(t.map((e=>{if(e in v)return;v[e]=!0;const t=e.endsWith(".css"),s=t?'[rel="stylesheet"]':"";if(document.querySelector(`link[href="${e}"]${s}`))return;const r=document.createElement("link");return r.rel=t?"stylesheet":h,t||(r.as="script",r.crossOrigin=""),r.href=e,document.head.appendChild(r),t?new Promise(((e,t)=>{r.addEventListener("load",e),r.addEventListener("error",t)})):void 0}))).then((()=>e()))},L=l({history:a(),routes:[{path:"/",redirect:"/maps"},{path:"/maps",component:()=>E((()=>__import__("./Maps.96f9045e.js")),["/assets/Maps.96f9045e.js","/assets/Maps.c86feea0.css","/assets/vendor.b62f66d0.js"])},{path:"/mods",component:()=>E((()=>__import__("./Mods.c51a30f9.js")),void 0)},{path:"/settings",component:()=>E((()=>__import__("./Settings.ee9c9663.js")),["/assets/Settings.ee9c9663.js","/assets/Settings.ba27059e.css","/assets/vendor.b62f66d0.js"])}]});c(_).use(L).mount("#app");
