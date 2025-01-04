# Video Overlayer

The automation to concatenate two videos for overlaying base video.

It's a part of a software suite for managing content in Telegram channels and VK groups. 

## Installation

To install the Video Overlayer application, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/sunchess/video_overlayer.git
    ```

2. Navigate to the project directory:
    ```sh
    cd video_overlayer
    ```

3. Install the dependencies:
    ```sh
    go mod download
    ```

4. Set up the environment variables. Create a [.env](http://_vscodecontentref_/0) file in the root directory of the project and add the necessary environment variables:
    ```env
    DB_PATH=your_database_path
    LOGS_PATH=your_logs_path
    POSTS_DIR=your_posts_directory
    SIGN_VIDEO_PATH=your_sign_video_path
    ```

5. Build the application:
    ```sh
    go build -o video_overlayer
    ```

6. Run the application:
    ```sh
    ./video_overlayer
    ```

Now the Video Overlayer application should be up and running.

## Features

- **Batch Processing**: All posts are divided into groups based on the number of posts specified in the configuration `GroupForProcessing` (`config/app.go`). This allows for efficient batch processing of posts.
- **Worker Distribution**: The posts are distributed among workers, with the number of workers also specified in the configuration `WorkerCount`. This ensures that the workload is balanced and processed in parallel, improving the overall efficiency of the application.
- **Video Concatenation**: The application selects a main video file from the `media` directory and randomly selects an overlay video from the `overlay` directory using the `SIGNS_VIDEO_PATH` environment variable. The selected overlay video is concatenated to the end of the main video.
- **Resolution Adjustment**: The resolution of the videos is adjusted to 720x1280 pixels (vertical video) to ensure consistency.
- **Output**: The resulting concatenated video is saved in the `media` directory with the name `processed.mp4`.

All processed videos are saved in the database, and this data can be used for posting to other social networks.

## Directory Structure

The directory structure for the Video Overlayer application should be as follows:

`posts_dir/{post_id}/media`


- `posts_dir`: The root directory containing all posts.
- `{post_id}`: A unique identifier for each post.
- `media`: A subdirectory within each post directory where the main video file (`.mp4`) is stored.

The application processes the videos as follows:
- It selects a main video file from the `media` directory.
- It randomly selects a video from the `overlay` directory uses `SIGN_VIDEO_PATH` env variable.
- The selected overlay video is concatenated to the end of the main video.
- The resolution of the videos is adjusted to 720x1280 pixels (vertical video) to ensure consistency.

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.