#!/usr/bin/env -S deno run -A lighthouse.ts

import lighthouse from "npm:lighthouse";
import * as chromeLauncher from "npm:chrome-launcher";

await Deno.mkdir("public/report/", { recursive: true });
const chrome = await chromeLauncher.launch({ chromeFlags: ["--headless"] });

const urls = [
	"/eu/ec/eci/fr.html",
	"/eu/ec/eci/2022/2/fr.html",
	"/eu/ec/eci/schema.html",

	"/fr.html",
	"/about/fr.html",
	"/release/",
	"/release/2024/",
];

let l = null;
const results = [];
for (const u of urls) {
	const { report, lhr } = await lighthouse(
		"http://localhost:8000" + u,
		{
			logLevel: "warn",
			locale: "fr",
			output: "html",
			port: chrome.port,
		},
	);
	l = lhr;

	const name = u.replaceAll("/", "_") + ".report.html";
	await Deno.writeTextFile(
		"public/report/" + name,
		report,
	);

	results.push({
		name,
		score: Object.entries(lhr.categories).reduce(
			(
				o,
				[categorie, { score }],
			) => (score != 1 && (o[categorie] = score * 100), o),
			{},
		),
		date: Intl.DateTimeFormat("fr", { timeStyle: "long" }).format(
			new Date(),
		),
	});
}

console.table(results);

Deno.writeTextFile(
	"public/report/index.html",
	`<!DOCTYPE html><html lang=en><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=0.5"><title>Lighthouse results</title></head><body><table style="font-size:large;font-family:mono">` +
		`<tr><th>Name<th>Score<th>Date` +
		results.map(({ name, score, date }) =>
			`<tr><td><a href="${name}">${name}</a><td>${
				JSON.stringify(score)
			}<td>${date}`
		).join(""),
);

chrome.kill();
