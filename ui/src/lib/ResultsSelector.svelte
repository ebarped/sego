<!-- https://www.skeleton.dev/utilities/popups#combobox FIRST INSTALL FLOATING-UI! -->
<script lang="ts">
    import { ListBox, ListBoxItem } from "@skeletonlabs/skeleton";
    import { popup } from "@skeletonlabs/skeleton";
    import type { PopupSettings } from "@skeletonlabs/skeleton";

    export let rscount: number;

    const popupCombobox: PopupSettings = {
        event: "hover-click", // Represents the type of event that opens/closed the popup
        target: "popupCombobox", // Matches the data-popup value on your popup element
        placement: "bottom", // Defines which side of your trigger the popup will appear
        closeQuery: ".listbox-item", // closes popup
    };
</script>

<button
    class="btn variant-filled w-48 justify-between"
    use:popup={popupCombobox}
>
    <span>Nº of Results: {rscount}</span>
    <span>↓</span>
</button>

<div class="card w-48 shadow-xl py-2" data-popup="popupCombobox">
    <ListBox rounded="rounded-none">
        <ListBoxItem bind:group={rscount} name="rscount" value="5">
            Nº of Results: 5
        </ListBoxItem>
        <ListBoxItem bind:group={rscount} name="rscount" value="10">
            Nº of Results: 10
        </ListBoxItem>
        <ListBoxItem bind:group={rscount} name="rscount" value="15">
            Nº of Results: 15
        </ListBoxItem>
    </ListBox>
    <div class="arrow bg-surface-100-800-token" />
</div>
