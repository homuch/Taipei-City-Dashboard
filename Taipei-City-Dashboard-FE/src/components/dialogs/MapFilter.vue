<!-- eslint-disable indent -->
<!-- eslint-disable no-mixed-spaces-and-tabs -->
<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { computed, watch, ref } from "vue";
import { useDialogStore } from "../../store/dialogStore";
import { useContentStore } from "../../store/contentStore";
import DialogContainer from "./DialogContainer.vue";

const timeout = ref(null);
const inputValue = ref("");
const msPerDay = 24 * 60 * 60 * 1000;

const dialogStore = useDialogStore();
const contentStore = useContentStore();
const filter_type = dialogStore.moreInfoContent.filter_type;
const canNotUseCoor = computed(() => {
	return (
		filter_type !== "time_clip" &&
		!(
			"geolocation" in navigator &&
			dialogStore.coor.latitude &&
			dialogStore.coor.longitude
		)
	);
});
const getPosition = () => {
	if ("geolocation" in navigator) {
		navigator.geolocation.getCurrentPosition((position) => {
			dialogStore.setCoor(
				position.coords.latitude,
				position.coords.longitude
			);
		});
		return true;
	} else {
		return false;
	}
};

const handleRadiusChange = () => {
	// console.log(e.target.value);
	// dialogStore.setMapFilterRadius(e.target.value);
	// console.log("Hellp");
	console.log(dialogStore.moreInfoContent);
	clearTimeout(timeout.value);
	timeout.value = setTimeout(() => {
		// dialogStore.mapFilterRadius = (
		// 	Date.now() - dialogStore.mapFilterRadius
		// ).getTime();
		// console.log(dialogStore.mapFilterRadius);
		console.log(filter_type);

		if (filter_type === "time_clip") {
			const dt = new Date();
			const corr = new Date(dt + inputValue.value * msPerDay);
			// console.log(corr);
			dialogStore.mapFilterRadius = `${corr.getFullYear()}-${
				corr.getMonth() + 1 < 10 ? "0" : ""
			}${corr.getMonth() + 1}-${
				corr.getDate() < 10 ? "0" : ""
			}${corr.getDate()} 00:00:00`;
		} else {
			dialogStore.mapFilterRadius = inputValue.value;
		}
		// console.log(dialogStore.mapFilterRadius);
		contentStore.getCurrentComponentData(dialogStore.moreInfoContent.index);
	}, 200);
};

watch(inputValue, handleRadiusChange);
watch(filter_type, () => {
	inputValue.value = 0;
});

const currentPosition = computed(() => {
	return dialogStore.coor;
});

const currentRepresentDate = computed(() => {
	let dt = new Date();
	let corr = new Date(dt - inputValue.value * msPerDay);
	return corr.toLocaleDateString();
});

function handleClose() {
	dialogStore.dialogs.mapFilter = false;
}
</script>

<template>
	<DialogContainer dialog="mapFilter" @on-close="handleClose">
		<div class="mapFilter">
			<h2>進階篩選資料</h2>
			<div class="mapFilter-control">
				<button
					v-if="filter_type === 'distance'"
					class="mapFilter-control-confirm"
					@click="getPosition"
				>
					使用現在位置資訊
				</button>

				<input
					v-if="filter_type === 'time_clip'"
					type="range"
					min="-30"
					max="30"
					step="1"
					v-model="inputValue"
				/>
				<input
					v-else
					type="range"
					min="0"
					max="15"
					step="0.1"
					v-bind:disabled="canNotUseCoor"
					v-model="inputValue"
				/>
			</div>
			<div class="mapFilter-input">
				<!-- 緯度：{{ currentPosition.latitude || "未知" }} 經度：{{
					currentPosition.longitude || "未知"
				}} -->
				<!-- 篩選半徑：{{ dialogStore.mapFilterRadius }} 公里 -->
				<!-- currentRepresentDate: {{ currentRepresentDate }} -->
				<span v-if="filter_type === 'time_clip'"
					>篩選日期：{{ currentRepresentDate }}</span
				>
				<span v-else>篩選半徑：{{ inputValue }} 公里</span>
			</div>
		</div>
	</DialogContainer>
</template>

<style scoped lang="scss">
.mapFilter {
	width: 300px;

	h3 {
		margin-bottom: 0.5rem;
		font-size: var(--font-s);
		font-weight: 400;
		color: var(--color-complement-text);
	}

	&-input {
		display: flex;
		flex-direction: column;
		margin: 0.5rem 0;

		textarea {
			font-size: var(--font-s);
			height: 160px;
			resize: none;
		}
	}

	&-control {
		display: flex;
		justify-content: flex-end;
		flex-direction: column;
		&-confirm {
			margin: 0 2px;
			padding: 4px 10px;
			border-radius: 5px;
			background-color: var(--color-highlight);
			transition: opacity 0.2s;

			&:hover {
				opacity: 0.8;
			}
		}
	}
}
</style>
