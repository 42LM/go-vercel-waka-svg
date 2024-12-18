# go-vercel-waka-svg
Create an SVG with dynamic content from [Wakatime](https://wakatime.com/) with a [Vercel](https://vercel.com/) serverless function that generates an SVG image containing the last 7 days of your WakaTime coding activity.

## Use in markdown
```
![]()
```
or
```
<img src="">
```

## Local development
First of all make sure the env variable `WAKA_API_KEY` is set.
```
cp .envrc.example .envrc
```

Run the app:
```
go run ./cmd/main.go
```

Example request:
```
curl localhost:8080/api?type=waka
```

## Deployment
Install vercel:
```
npm i -g vercel
```

Deploy the current directory:
```
vercel
```

Deploy to production:
```
vercel --prod
```
