<template>
  <v-container>
    <v-card class="mx-auto" max-width="400" hover>
      <v-flex>
        <v-img contain max-height="100px" :src="data.image"></v-img>
      </v-flex>
      <v-card-title class="justify-center">
        <div>{{ data.metadata.name }}</div>
      </v-card-title>
      <!-- <v-card-text>
        <v-chip v-if="data.metadata.labels.chip" x-small class="justify-center">
          {{ data.metadata.labels.chip }}</v-chip
        >
      </v-card-text> -->
      <v-divider></v-divider>
      <v-card-text class="text--primary">
        <div class="mb-4">
          {{ data.metadata.selfLink }}
        </div>
      </v-card-text>
      <v-spacer></v-spacer>
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