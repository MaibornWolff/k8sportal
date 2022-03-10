<template>
  <v-layout class="ma-3">
    <v-flex xs12>
      <v-sheet>
        <v-container>
          <v-row justify="center">
            <v-col>
              <h1>Services</h1>
            </v-col>

            <v-col md="2">
              <v-text-field
                label="Search"
                prepend-icon="mdi-magnify"
                v-on:input="resultQuery($event)"
              ></v-text-field>
            </v-col>
            <v-col md="3" offset-md="2" class="d-flex align-center">
              <h3>Filter:</h3>
              <v-chip-group active-class="primary--text" column>
                <v-chip
                  v-for="(chip, key) in chips"
                  :key="key"
                  x-small
                  class="ml-1"
                  @click="filterTiles(chip)"
                  >{{ chip }}</v-chip
                >
              </v-chip-group>
            </v-col>
          </v-row>

          <v-row>
            <v-flex v-for="(tile, key) in renderTile" :key="key">
              <v-card flat tile>
                <Tile :data="tile" />
              </v-card>
            </v-flex>
          </v-row>
        </v-container>
      </v-sheet>
    </v-flex>
  </v-layout>
</template>

<script>
import Tile from "./Tile.vue";
import axios from "axios";

export default {
  name: "GridList",
  components: {
    Tile,
  },
  data() {
    return {
      tiles: [],
      images: [],
      names: [],
      chips: ["keiner"],
      selectedTiles: "all",
      filteredTiles: [],
      searchResult: [],
      isSearch: false,
    };
  },

  methods: {
    // imports all image path in assets/ folder -> push in images
    importAll(r) {
      r.keys().forEach((key) =>
        this.images.push({ pathLong: r(key), pathShort: key })
      );
    },
    // get short name of all image paths -> push in names
    getLogoNames() {
      const re = /(?<=.\/)(.*?)(?=.svg)/g;
      this.importAll(require.context("../assets/", true, /\.svg$/));
      this.images.map((element) =>
        this.names.push(element.pathShort.match(re)[0])
      );
    },
    // for each tile and for each name -> if match -> insert new image attribute in tile object with image path of matched image
    // and pass via props
    getImages() {
      const img = require.context("../assets/", false, /\.svg$/);
      this.tiles.map((tile) => {
        this.names.forEach((name) => {
          if (!tile.metadata.name.includes(name)) {
            tile["image"] = img("./" + "default" + ".svg");
          }
        });

        this.names.forEach((name) => {
          if (tile.metadata.name.includes(name)) {
            tile["image"] = img("./" + name + ".svg");
          }
        });
      });
    },
    getAvailableChips() {
      this.tiles.map((tile) => {
        if (!this.chips.includes(tile.metadata.labels.chips)) {
          this.chips.push(tile.metadata.labels.chips);
        }
      });
    },
    filterTiles(chip) {
      this.isActive = !this.isActive;
      if (chip == "keiner") {
        this.selectedTiles = "all";
      } else {
        this.filteredTiles = this.tiles.filter((tile) => {
          return tile.metadata.labels.chips == chip;
        });
        this.selectedTiles = "filter";
      }
    },
    resultQuery(event) {
      const searchQuery = event;
      if (searchQuery) {
        this.isSearch = true;
        this.searchResult = this.tiles.filter((tile) => {
          if (tile.metadata.labels.serviceName) {
            return searchQuery
              .toLowerCase()
              .split(" ")
              .every((v) =>
                tile.metadata.labels.serviceName.toLowerCase().includes(v)
              );
          } else {
            return searchQuery
              .toLowerCase()
              .split(" ")
              .every((v) => tile.metadata.name.toLowerCase().includes(v));
          }
        });
      } else {
        this.isSearch = false;
      }
    },
  },
  mounted() {
    axios
      .get(this.api_endpoint)
      .then((res) => {
        this.tiles = res.data;

        this.getLogoNames();
        this.getImages();
        this.getAvailableChips();
      })
      .catch((err) => console.log(err));
  },
  computed: {
    // Dynamic env variable
    api_endpoint: function() {return this.VUE_APP_API_ENDPOINT},
    // filter by tag
    renderTile() {
      return this[this.selectedTiles];
    },
    all() {
      if (this.isSearch) {
        return this.searchResult;
      } else {
        return this.tiles;
      }
    },
    filter() {
      if (this.isSearch) {
        return this.searchResult.filter((value) =>
          this.filteredTiles.includes(value)
        );
      } else {
        return this.filteredTiles;
      }
    },
  },
};
</script>
