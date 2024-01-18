// Packages
import fs from "fs";
import path from "path";
import * as database from "./v4-database/prisma.js";
import * as perms from "./perms.js";
import cors from "@fastify/cors";
import ratelimit from "@fastify/rate-limit";
import "dotenv/config";
import Fastify, {
	FastifyInstance,
	FastifyReply,
	FastifyRequest,
} from "fastify";

// Middleware
const app: FastifyInstance = Fastify({
	logger: true,
});
let Routes = [];

app.register(cors, {
	origin: "*",
	allowedHeaders: [
		"secret",
		"userid",
		"Authorization",
		"Authorization",
		"Content-Type",
		"Content-Disposition",
		"Content-Length",
	],
	methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
	credentials: true,
	optionsSuccessStatus: 200,
	preflight: true,
	strictPreflight: false,
});

app.register(ratelimit, {
	global: true,
	max: 50,
	timeWindow: 1000,
});

app.addHook("preHandler", (req, res, done) => {
	res.header("Access-Control-Allow-Origin", "*");
	res.header("Access-Control-Allow-Headers", "*");
	res.header("Access-Control-Allow-Methods", "*");
	res.header("Access-Control-Allow-Credentials", "true");

	done();
});

// Endpoint for getting Routes
app.route({
	method: "GET",
	url: "/routes",
	handler: async (request: FastifyRequest, reply: FastifyReply) => {
		return reply.send(Routes);
	},
});

// API Endpoints Map
const getFilesInDirectory = (dir: string) => {
	let files: string[] = [];
	const filesInDir = fs.readdirSync(dir);

	for (const file of filesInDir) {
		const filePath = path.join(dir, file);
		const stat = fs.statSync(filePath);

		if (stat.isDirectory())
			files = files.concat(getFilesInDirectory(filePath));
		else files.push(filePath);
	}

	return files;
};

// API Endpoints
const apiEndpointsFiles = getFilesInDirectory("./dist/endpoints").filter(
	(file) => file.endsWith(".js")
);

for (const file of apiEndpointsFiles) {
	import(`../${file}`)
		.then((module) => {
			app.route(module.default);
			Routes.push(module.default);
		})
		.catch((error) => {
			console.error(`Error importing ${file}: ${error}`);
		});
}

setTimeout(() => {
	// Start Server
	app.listen({ port: Number(process.env.PORT) }, (err) => {
		if (err) throw err;
	});
}, 5000);