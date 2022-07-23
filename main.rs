use async_std::fs::File;
use async_std::io::BufReader;
use async_std::path::Path;
use async_std::prelude::*;
use async_std::sync::Mutex;
use async_std::{fs, task};
use std::sync::Arc;
use std::time::Instant;

async fn add_word(word: String, words: Arc<Mutex<Vec<String>>>) {
    println!("{}", word);
    let mut unlocked = words.lock().await;
    unlocked.push(word);
    drop(unlocked);
}

fn main() {
    let start = Instant::now();
    let mut word_count = 0;
    let words: Arc<Mutex<Vec<String>>> = Arc::new(Mutex::new(Vec::new()));

    let path = Path::new(".")
        .join("data")
        .join("a");
    let mut selected_word_count: usize = 0;

    task::block_on(async {
        let mut entries = fs::read_dir(path).await.expect("Unable to list");
        while let Some(res) = entries.next().await {
            word_count += 1;
            let entry = res.expect("Unable to get entry");
            //println!("{}", entry.file_name().to_string_lossy());
            let words_clone = words.clone();

            task::spawn(async move {
                let frequency = entry.path().join("2015-2017-spoken-frequency.txt");
                let file = File::open(&frequency)
                    .await
                    .expect("Unable to open frequency");
                let mut lines = BufReader::new(file).lines();
                while let Some(line) = lines.next().await {
                    if 798 <= line.expect("Unable to read line").parse::<u32>().unwrap() {
                        //println!("{}", entry.file_name().to_os_string().into_string().unwrap());
                        let frequency_clone = frequency.clone();

                        task::spawn(async move {
                            let mut word = String::from("");
                            let translation = frequency_clone
                                .parent()
                                .expect("Unable to unable")
                                .join("translation.txt");
                            let file = File::open(&translation)
                                .await
                                .expect("Unable to open translation");

                            let mut lines = BufReader::new(file).lines();
                            while let Some(line) = lines.next().await {
                                word = line.expect("Unable to read line");
                                break;
                            }

                            if word.is_empty() {
                                word = translation
                                    .parent()
                                    .expect("Unable to unable")
                                    .file_name()
                                    .expect("Unable to unable")
                                    .to_os_string()
                                    .into_string()
                                    .unwrap();
                            }
                            add_word(word, words_clone).await;
                        });
                    }
                    break;
                }
            });
        }
        selected_word_count = words.lock().await.len();
    });
    println!(
        "selected {} words from {}, took {} seconds",
        selected_word_count,
        word_count,
        start.elapsed().as_millis() as f64 / 1000.0
    );
}
