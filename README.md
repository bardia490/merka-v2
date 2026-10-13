# Goals

the goal for this project is to rewrite the previous merka app (written in D programming language)
in Go. the rewrite will be 100% backwards compatible with the features of the older version. when the rewrite is done
the app will also be getting new updates to its UI alongside new features and performance improvemnts.

# PLANED FOR THE FUTURE

- Add the option Write work information to file
- Add discount mechanism for printing work prices
- Change all fmt.Print* functions to os.Stdout.Write (or something similar) for better performance and GC performance
- Make the print formating better
- Add the option to edit work, monjog, price_tags and ...
- Add better error formating and printing
- Clean up and add better UI for the main file in cmd
- Fix all TODO's
