# Take home test
A go service to retrieve repository information and commits from github, save the data in a persistent store(postgres),and continuously monitor the repository for changes

### environment variables
```
DATABASE_URL=postgres://user:password@localhost:5455/github_db?search_path=github_service
OWNER=chromium
REPO=chromium
```
An optional ```SINCE``` env variable is used for configurable date to start pulling commits since and should be in format YYYY-MM-DDTHH:MM:SSZ

### HOW TO RUN PROJECT MANUALLY
To run project, make sure the available environment variables listed above are present, then:
```
1. build code with command go build -o bin/github-service ./main.go
2. If you want to manually run the script you can run ./bin/github-service cron-monitor-repo
```

### CONTINUOUSLY MONITOR REPO
```
To continuously monitor repo and run this service you can use cronjob with example shown in cron.example
Replace the path with correct path to file.
```

### DOCUMENTATION
*To get the top N commit authors by commit counts from the database*;

Select the *author* and count the number of commits for each author, using *COUNT(*)*. The result is aliased as *commit_count*.
*Group* the results by the author column, so the COUNT(*) function counts the number of commits for each distinct author.
*Orders* the results by the commit_count in descending order, so the authors with the most commits appear first.
*Limits* the number of results to the top N authors, where N is the number of authors you want to retrieve. An example to count top 5 committers by commit counts 
```
SELECT author, COUNT(*) as commit_count
FROM commit
GROUP BY author
ORDER BY commit_count DESC
LIMIT 5;
```

*To retrieve commits of a repository by repository name from the database*;
Since we have a reference between the repository table and the commit table using the unique repo_full_name field 
we can easily retrieve all commits for a repository using a join statment where the commit.repo_full_name = repository.repo_full_name. example is shown below
```
SELECT c.id, c.sha, c.message, c.url, c.author, c.date FROM github_service.commit c
JOIN github_service.repository r ON c.repo_full_name = r.repo_full_name
where r.repo_full_name='chromium/chromium'
```


