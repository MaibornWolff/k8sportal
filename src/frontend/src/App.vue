<template>
  <v-app>
    <v-app-bar app color="light-grey" prominent>
      <div class="d-flex align-center">
        <v-img
          alt="MaibornWolff Logo"
          src="./assets/maibornwolff.svg"
          transition="scale-transition"
          max-width="150"
          class="ma-auto mx-4"
          contain
        />
      </div>
      <v-spacer />
      <v-toolbar-title class="ma-auto mx-4 text-h4 text-right">
        Cluster Portal
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

  data: () => ({
    //
  }),
};
</script>
