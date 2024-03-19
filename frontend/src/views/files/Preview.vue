<template>
  <div
    id="previewer"
    @mousemove="toggleNavigation"
    @touchstart="toggleNavigation"
  >
    <header-bar @update:hideNotifications="hideNotifications">
      <action icon="close" :label="$t('buttons.close')" @action="close()" />
      <title>{{ name }}</title>
      <action
        :disabled="loading"
        v-if="isResizeEnabled && req.type === 'image'"
        :icon="fullSize ? 'photo_size_select_large' : 'hd'"
        @action="toggleSize"
      />

      <template #actions>
        <action
          :disabled="loading"
          v-if="user.perm.rename"
          icon="mode_edit"
          :label="$t('buttons.rename')"
          show="rename"
        />
        <action
          :disabled="loading"
          v-if="user.perm.delete"
          icon="delete"
          :label="$t('buttons.delete')"
          @action="deleteFile"
          id="delete-button"
        />
        <action
          :disabled="loading"
          v-if="user.perm.download"
          icon="file_download"
          :label="$t('buttons.download')"
          @action="download"
        />
        <action
          :disabled="loading"
          icon="info"
          :label="$t('buttons.info')"
          show="info"
        />
        <action
          icon="notifications"
          :label="$t('buttons.notifications')"
          :counter="notificationIconCount"
          @action="openNotifications"
          class="notifications"
        />
      </template>
    </header-bar>

    <div class="loading delayed" v-if="loading">
      <div class="spinner">
        <div class="bounce1"></div>
        <div class="bounce2"></div>
        <div class="bounce3"></div>
      </div>
    </div>
    <template v-else>
      <div
        class="preview"
        v-bind:class="{
          previewShowComments: showComments,
          previewHideComments: !showComments,
        }"
      >
        <ExtendedImage v-if="req.type == 'image'" :src="raw"></ExtendedImage>
        <audio
          v-else-if="req.type == 'audio'"
          ref="player"
          :src="raw"
          controls
          :autoplay="autoPlay"
          @play="autoPlay = true"
        ></audio>
        <video
          v-else-if="req.type == 'video'"
          ref="player"
          :src="raw"
          controls
          :autoplay="autoPlay"
          @play="autoPlay = true"
        >
          <track
            kind="captions"
            v-for="(sub, index) in subtitles"
            :key="index"
            :src="sub"
            :label="'Subtitle ' + index"
            :default="index === 0"
          />
          Sorry, your browser doesn't support embedded videos, but don't worry,
          you can <a :href="downloadUrl">download it</a>
          and watch it with your favorite video player!
        </video>
        <object
          v-else-if="req.extension.toLowerCase() == '.pdf'"
          class="pdf"
          :data="raw"
        ></object>
        <div v-else-if="req.type == 'blob'" class="info">
          <div class="title">
            <i class="material-icons">feedback</i>
            {{ $t("files.noPreview") }}
          </div>
          <div>
            <a target="_blank" :href="downloadUrl" class="button button--flat">
              <div>
                <i class="material-icons">file_download</i
                >{{ $t("buttons.download") }}
              </div>
            </a>
            <a
              target="_blank"
              :href="raw"
              class="button button--flat"
              v-if="!req.isDir"
            >
              <div>
                <i class="material-icons">open_in_new</i
                >{{ $t("buttons.openFile") }}
              </div>
            </a>
          </div>
        </div>
      </div>
      <div
        class="notifications-pane"
        v-bind:class="{ frameShowNotifications: showNotifications }"
      >
        <div
          v-if="isMobile"
          class="notifications-mobile-bar notifications-mobile-bar-preview"
        >
          <action
            icon="close"
            :label="$t('buttons.close')"
            @action="openNotifications"
          />
        </div>

        <NotificationPane
          v-bind:filePath="req.path"
          @update:unacknowledgedNotificationCount="
            updateUnacknowledgedNotificationCount
          "
        >
        </NotificationPane>
      </div>
      <div class="socialContainer">
        <Social
          v-bind:filePath="req.path"
          v-bind:showComments="showComments"
          v-bind:screenWidth="width"
          @update:showComments="updateShowComments"
          @update:commentEntering="commentEntering"
        ></Social>
      </div>
    </template>

    <button
      @click="prev"
      @mouseover="hoverNav = true"
      @mouseleave="hoverNav = false"
      :class="{ hidden: !hasPrevious || !showNav }"
      :aria-label="$t('buttons.previous')"
      :title="$t('buttons.previous')"
    >
      <i class="material-icons noselect">chevron_left</i>
    </button>
    <button
      @click="next"
      @mouseover="hoverNav = true"
      @mouseleave="hoverNav = false"
      :class="{ hidden: !hasNext || !showNav }"
      :aria-label="$t('buttons.next')"
      :title="$t('buttons.next')"
    >
      <i class="material-icons noselect">chevron_right</i>
    </button>
    <link rel="prefetch" :href="previousRaw" />
    <link rel="prefetch" :href="nextRaw" />
  </div>
</template>

<script>
import { mapGetters, mapState } from "vuex";
import { files as api } from "@/api";
import { resizePreview } from "@/utils/constants";
import url from "@/utils/url";
import throttle from "lodash.throttle";
import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import ExtendedImage from "@/components/files/ExtendedImage.vue";
import Social from "@/components/files/Social.vue";
import NotificationPane from "@/components/NotificationPane.vue";

const mediaTypes = ["image", "video", "audio", "blob"];

export default {
  name: "preview",
  components: {
    HeaderBar,
    Action,
    ExtendedImage,
    Social,
    NotificationPane,
  },
  data: function () {
    return {
      previousLink: "",
      nextLink: "",
      listing: null,
      name: "",
      fullSize: false,
      showNav: true,
      navTimeout: null,
      hoverNav: false,
      autoPlay: false,
      previousRaw: "",
      nextRaw: "",
      showComments: true,
      showNotifications: false,
      unacknowledgedNotificationCount: 0,
      width: window.innerWidth,
      isCommentEntering: false,
    };
  },
  computed: {
    ...mapState(["req", "user", "oldReq", "jwt", "loading"]),
    ...mapGetters(["currentPrompt"]),
    hasPrevious() {
      return this.previousLink !== "";
    },
    hasNext() {
      return this.nextLink !== "";
    },
    downloadUrl() {
      return api.getDownloadURL(this.req);
    },
    raw() {
      if (this.req.type === "image" && !this.fullSize) {
        return api.getPreviewURL(this.req, "big");
      }

      return api.getDownloadURL(this.req, true);
    },
    showMore() {
      return this.currentPrompt?.prompt === "more";
    },
    isResizeEnabled() {
      return resizePreview;
    },
    subtitles() {
      if (this.req.subtitles) {
        return api.getSubtitlesURL(this.req);
      }
      return [];
    },
    notificationIconCount() {
      return this.unacknowledgedNotificationCount;
    },
    isMobile() {
      return this.width <= 736;
    },
  },
  watch: {
    $route: function () {
      this.updatePreview();
      this.toggleNavigation();
    },
  },
  async mounted() {
    window.addEventListener("keydown", this.key);
    window.addEventListener("resize", this.windowsResize);
    this.listing = this.oldReq.items;
    this.updatePreview();
  },
  beforeDestroy() {
    window.removeEventListener("keydown", this.key);
    window.removeEventListener("resize", this.windowsResize);
  },
  methods: {
    updateShowComments(newValue) {
      this.showComments = newValue;
    },
    commentEntering(newValue) {
      this.isCommentEntering = newValue;
    },
    deleteFile() {
      this.$store.commit("showHover", {
        prompt: "delete",
        confirm: () => {
          this.listing = this.listing.filter((item) => item.name !== this.name);

          if (this.hasNext) {
            this.next();
          } else if (!this.hasPrevious && !this.hasNext) {
            this.close();
          } else {
            this.prev();
          }
        },
      });
    },
    prev() {
      this.hoverNav = false;
      this.$router.replace({ path: this.previousLink });
    },
    next() {
      this.hoverNav = false;
      this.$router.replace({ path: this.nextLink });
    },
    key(event) {
      if (this.currentPrompt !== null) {
        return;
      }

      if (!this.isCommentEntering) {
        if (event.which === 39) {
          // right arrow
          if (this.hasNext) this.next();
        } else if (event.which === 37) {
          // left arrow
          if (this.hasPrevious) this.prev();
        } else if (event.which === 27) {
          // esc
          this.close();
        }
      }
    },
    async updatePreview() {
      if (
        this.$refs.player &&
        this.$refs.player.paused &&
        !this.$refs.player.ended
      ) {
        this.autoPlay = false;
      }

      let dirs = this.$route.fullPath.split("/");
      this.name = decodeURIComponent(dirs[dirs.length - 1]);

      if (!this.listing) {
        try {
          const path = url.removeLastDir(this.$route.path);
          const res = await api.fetch(path);
          this.listing = res.items;
        } catch (e) {
          this.$showError(e);
        }
      }

      this.previousLink = "";
      this.nextLink = "";

      for (let i = 0; i < this.listing.length; i++) {
        if (this.listing[i].name !== this.name) {
          continue;
        }

        for (let j = i - 1; j >= 0; j--) {
          if (mediaTypes.includes(this.listing[j].type)) {
            this.previousLink = this.listing[j].url;
            this.previousRaw = this.prefetchUrl(this.listing[j]);
            break;
          }
        }
        for (let j = i + 1; j < this.listing.length; j++) {
          if (mediaTypes.includes(this.listing[j].type)) {
            this.nextLink = this.listing[j].url;
            this.nextRaw = this.prefetchUrl(this.listing[j]);
            break;
          }
        }

        return;
      }
    },
    prefetchUrl(item) {
      if (item.type !== "image") {
        return "";
      }

      return this.fullSize
        ? api.getDownloadURL(item, true)
        : api.getPreviewURL(item, "big");
    },
    openMore() {
      this.$store.commit("showHover", "more");
    },
    resetPrompts() {
      this.$store.commit("closeHovers");
    },
    toggleSize() {
      this.fullSize = !this.fullSize;
    },
    toggleNavigation: throttle(function () {
      this.showNav = true;

      if (this.navTimeout) {
        clearTimeout(this.navTimeout);
      }

      this.navTimeout = setTimeout(() => {
        this.showNav = false || this.hoverNav;
        this.navTimeout = null;
      }, 1500);
    }, 500),
    close() {
      this.$store.commit("updateRequest", {});

      let uri = url.removeLastDir(this.$route.path) + "/";
      this.$router.push({ path: uri });
    },
    download() {
      window.open(this.downloadUrl);
    },
    openNotifications: function () {
      this.$store.commit("closeHovers");
      this.showNotifications = !this.showNotifications;
    },
    hideNotifications: function () {
      this.showNotifications = false;
    },
    updateUnacknowledgedNotificationCount: function (newValue) {
      this.unacknowledgedNotificationCount = newValue;
    },
    windowsResize: throttle(function () {
      this.width = window.innerWidth;
    }, 100),
  },
};
</script>
<style>
.notifications-mobile-bar-preview {
  padding-left: 1rem;
}

#previewer .previewHideComments {
  text-align: center;

  height: calc(100vh - 16em);
}

#previewer .previewShowComments {
  text-align: center;

  height: calc(50vh - 8em);
}

#previewer .socialContainer {
  background: var(--socialContainerBg);
  color: var(--textPrimary);
  overflow-y: auto;

  height: calc(50vh + 4em);
  padding-left: 2em;
  padding-right: 2em;
}

@media (min-width: 1400px) {
  #previewer .previewHideComments {
    height: calc(100vh - 4em);

    width: calc(100vw - 36em);
  }

  #previewer .previewShowComments {
    height: calc(100vh - 4em);

    width: calc(100vw - 36em);
  }

  #previewer .socialContainer {
    height: calc(100vh - 4em);

    position: fixed;
    top: 4em;
    right: 0;
    width: 36em;
    padding-left: 0em;
    padding-right: 0em;
  }
}

#previewer .showCommentsPrompt {
  text-align: center;
}

#previewer .hideCommentsPrompt {
  text-align: center;
}
</style>
