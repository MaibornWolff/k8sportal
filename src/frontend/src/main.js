import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import vuetify from "./plugins/vuetify";

Vue.config.productionTip = false;

const getDynamicEnv = async () => {
  const dynamicEnv = await fetch('/dynamicEnv.json');
  return await dynamicEnv.json()
}

getDynamicEnv().then(function(json) {
  console.debug("dynamic Env: " + JSON.stringify(json))
  Vue.mixin({
    data() {
      return json;
    },
  });

  new Vue({
    router,
    vuetify,
    render: (h) => h(App),
  }).$mount("#app");
});
