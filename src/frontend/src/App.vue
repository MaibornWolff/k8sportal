<template>
  <v-app>
    <v-app-bar app color="light-grey" prominent hide-on-scroll>
      <!-- <div class="d-flex align-center"></div> -->
      <v-img
        v-bind:src="require('./assets/' + customLogo)"
        transition="scale-transition"
        class="ma-auto"
        max-height="90"
        max-width="90"
        contain
      />
      <v-spacer />
      <v-toolbar-title class="ma-auto mx-4 text-h4 text-right">
        {{ customTitle }}
      </v-toolbar-title>
    </v-app-bar>

    <v-main class="blue-grey lighten-4">
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: "App",
  // Use computed for variables from the dynamic environment to create local
  // copies and avoid changes propagating through the entire
  // application.
  computed:  {
      customTitle: function() {return this.VUE_APP_CUSTOM_TITLE},
      customLogo: function() {return this.VUE_APP_CUSTOM_LOGO}
  },
  created() {
    document.title = this.customTitle;
  },
};
</script>
