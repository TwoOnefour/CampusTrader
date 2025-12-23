import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
// import Dashboard from "./components/Dashboard.vue"
//
// const app = createApp(App);
// app.component('MyGlobalComponent', Dashboard);
import naive from "naive-ui";
const app = createApp(App);

app.use(naive);

app.mount("#app");
