<template>
  <v-container>
    <v-card class="mx-auto" max-width="400" hover>
      <v-flex>
        <v-img contain max-height="100px" :src="data.image"> </v-img>
      </v-flex>
      <v-card-title class="justify-center">
        <div v-if="data.metadata.labels.serviceName">
          {{ data.metadata.labels.serviceName }}
        </div>
        <div v-else>{{ data.metadata.name }}</div>
      </v-card-title>

      <v-divider></v-divider>
      <v-card-text class="text--primary">
        <div class="mb-4">
          {{ data.metadata.selfLink }}
        </div>
      </v-card-text>

      <v-spacer></v-spacer>

      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn icon @click="show = !show">
          <v-icon>{{ show ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
        </v-btn>
      </v-card-actions>

      <v-expand-transition>
        <div v-show="show">
          <v-divider></v-divider>
          <v-card-text v-if="data.metadata.annotations.description">
            {{ data.metadata.annotations.description }}
          </v-card-text>
          <v-card-text v-else> No description available. </v-card-text>
        </div>
      </v-expand-transition>
    </v-card>
  </v-container>
</template>

<script>
export default {
  props: ["data"],
  data() {
    return {
      show: false,
    };
  },
};
</script>