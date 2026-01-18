<script lang="ts">
	import InfiniteLoading from 'svelte-infinite-loading';
	
	const api = 'https://hn.algolia.com/api/v1/search_by_date?tags=story';
	
	type HackerNewsItem = {
		objectID: string;
		title: string;
		url: string | null;
		points: number;
		author: string;
		num_comments: number;
	};
	
	type InfiniteHandlerDetail = {
		loaded: () => void;
		complete: () => void;
	};
	
	type ApiResponse = {
		hits: HackerNewsItem[];
	};
	
	let page = 1;
	let list: HackerNewsItem[] = [];
	
	function infiniteHandler({ detail: { loaded, complete } }: CustomEvent<InfiniteHandlerDetail>) {
		console.log('list', list);

		fetch(`${api}&page=${page}`)
			.then(response => response.json() as Promise<ApiResponse>)
			.then(data => {
				if (data.hits.length) {
					page += 1;
					list = [...list, ...data.hits];
					loaded();
				} else {
					complete();
				}
			});
	}
</script>

<div id="app">
  
	{#each list as item, index}
		<div class="hacker-news-item" data-num={index + 1}>
			<a target="_blank" href={item.url}>{item.title}</a>
			<p>
				<span>{item.points}</span>
				points by
				<a target="_blank" href="https://news.ycombinator.com/user?id={item.author}">{item.author}</a>
				|
				<a target="_blank" href="https://news.ycombinator.com/item?id={item.objectID}">{item.num_comments} comments</a>
			</p>
		</div>
	{/each}
	
	<InfiniteLoading on:infinite={infiniteHandler}>
	   <div
		slot="spinner"
		class="flex justify-center w-full py-4">
		 <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24"><circle cx="18" cy="12" r="0" fill="currentColor"><animate attributeName="r" begin=".67" calcMode="spline" dur="1.5s" keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8" repeatCount="indefinite" values="0;2;0;0"/></circle><circle cx="12" cy="12" r="0" fill="currentColor"><animate attributeName="r" begin=".33" calcMode="spline" dur="1s" keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8" repeatCount="indefinite" values="0;2;0;0"/></circle><circle cx="6" cy="12" r="0" fill="currentColor"><animate attributeName="r" begin="0" calcMode="spline" dur="1.5s" keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8" repeatCount="indefinite" values="0;2;0;0"/></circle></svg>
	   </div>
    </InfiniteLoading>

</div>

<style>
 :global(body) {
	 padding-top: 28px;
	 background-color: #f6f6ef;
 }

 .hacker-news-item {
	 margin: 10px 0;
	 padding: 0 10px 0 40px;
	 line-height: 16px;
	 font-size: 14px;
 }
 .hacker-news-item::before {
	 content: attr(data-num) '.';
	 float: left;
	 margin-left: -40px;
	 width: 32px;
	 color: #888;
	 text-align: right;
 }
 .hacker-news-item > a {
	 color: #333;
	 text-decoration: none;
 }
 .hacker-news-item > a:hover {
	 color: #000;
 }
 .hacker-news-item p {
	 margin: 0;
	 font-size: 12px;
 }
 .hacker-news-item p, .hacker-news-item p a {
	 color: #888;
 }
 .hacker-news-item p a:not(:hover) {
	 text-decoration: none;
 }
</style>
