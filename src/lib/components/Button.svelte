<script lang="ts">
    import { cn } from "$lib/utils";

    let className: string = "";
    export { className as class };
</script>

<button
    class={cn(
        "inline-flex shimmer-button items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2",
        className,
    )}
    {...$$restProps}
>
    <slot />
</button>