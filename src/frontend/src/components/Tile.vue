<template>
  <v-container>
    <v-card class="mx-auto" width="300" hover>
      <v-list-item one-line>
        <v-list-item-content>
          <v-list-title>
            <v-chip x-small>
              {{ data.category }}
            </v-chip>
          </v-list-title>
        </v-list-item-content>
      </v-list-item>
      <v-flex>
        <v-img contain max-height="100px" :src="data.image"> </v-img>
      </v-flex>
      <v-card-title class="justify-center">
        <div>
          {{ data.serviceName }}
        </div>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-actions>
        <v-btn
          color="blue-grey"
          class="ma-2 white--text"
          :href="'http://'+data.ingressRules[0].ingressHost"
        >
          Service
          <v-icon right dark> mdi-chevron-right </v-icon>
        </v-btn>
        <v-spacer></v-spacer>

        <v-btn icon @click="show = !show">
          <v-icon>{{ show ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
        </v-btn>
      </v-card-actions>

      <v-expand-transition>
        <div v-show="show">
          <v-divider></v-divider>
          <v-card-text v-if="data.description">
            {{ data.description }}
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
