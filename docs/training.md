# Recommended Fine-tuning Settings

This document outlines the recommended parameters and configurations for fine-tuning the large language model in APISeeker.

> The following parameters were used in the experiments for our paper and can be adjusted as needed.

## Model Loading & Quantization

To load and train large models with limited GPU memory, we have adopted the following optimization measures:

- **Quantization**:
    - **`load_in_8bit=True`**: Enables loading the model with 8-bit quantization.

- **Data Type**:
    - **`torch_dtype=torch.bfloat16`**: Uses the `bfloat16` data type for model computations.

- **Attention Mechanism**:
    - **`attn_implementation="flash_attention_2"`**: Employs the Flash Attention 2 mechanism.

## PEFT & LoRA Settings

We use the LoRA (Low-Rank Adaptation) method from the PEFT (Parameter-Efficient Fine-Tuning) library for efficient fine-tuning.

## Training Parameters (SFTConfig)

- **Batch Size**:
    - **`per_device_train_batch_size=1`**: Sets the training batch size per device to 1.
    - **`gradient_accumulation_steps=8`**: Sets the number of gradient accumulation steps to 8. This means gradients from 8 batches (`1 * 8`) are accumulated before a parameter update is performed.

- **Performance & Optimization**:
    - **`bf16=True`**: Enables `bfloat16` mixed-precision training.
    - **`optim="adamw_apex_fused"`**: Uses the fused AdamW optimizer from the NVIDIA Apex library.
    - **`gradient_checkpointing=True`**: Enables gradient checkpointing.

- **Learning Rate & Scheduler**:
    - **`learning_rate=1.0e-4`**: Sets the learning rate to 0.0001.
    - **`lr_scheduler_type="linear"`**: Uses a linear learning rate scheduler.
    - **`warmup_ratio=0.1`**: Uses the first 10% of training steps for warmup, during which the learning rate gradually increases from 0 to the set value.

- **Regularization**:
    - **`weight_decay=0.0`**: Sets the weight decay to 0.

- **Training & Logging**:
    - **`num_train_epochs=1`**: Sets the number of training epochs to 1, meaning the dataset is traversed completely once.
    - **`max_seq_length=2048`**: Sets the maximum input sequence length to 2048 tokens.
    - **`logging_steps=10`**: Logs training information every 10 steps.
